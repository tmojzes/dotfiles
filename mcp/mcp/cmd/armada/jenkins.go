package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	jenkinsBase    = "https://alchemy-containers-jenkins.swg-devops.com"
	jenkinsJobPath = "/job/Containers-Runtime/job/armada-deploy-get-master-info"
)

var jenkinsJobURL = jenkinsBase + jenkinsJobPath

type jenkinsClient struct {
	httpClient *http.Client
	user       string
	token      string
	cfg        serverConfig
}

// newJenkinsClient reads the hand-exported Jenkins credentials from the
// environment. Missing credentials are a per-call tool error, not a startup
// failure, so the server stays connectable without them.
func newJenkinsClient(cfg serverConfig) (*jenkinsClient, error) {
	user := os.Getenv("JENKINS_USER")
	token := os.Getenv("JENKINS_API_KEY")
	if user == "" || token == "" {
		return nil, errors.New("JENKINS_USER and JENKINS_API_KEY must be set (token from ~/keys/jenkins-api-key)")
	}
	return &jenkinsClient{
		httpClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// The trigger POST must not follow the queue redirect;
				// the Location header is the payload we need.
				if req.Method == http.MethodPost {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
		user:  user,
		token: token,
		cfg:   cfg,
	}, nil
}

func (c *jenkinsClient) do(ctx context.Context, method, rawURL string, form url.Values) (*http.Response, error) {
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.user, c.token)
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return c.httpClient.Do(req)
}

func (c *jenkinsClient) triggerJob(ctx context.Context, clusterID string, commands []string) (string, error) {
	form := url.Values{}
	form.Set("CUSTOM_KUBX_KUBECTL", strings.Join(commands, "; "))
	form.Set("CLUSTER", clusterID)
	form.Set("CLUSTER_ETCD_INFO", "false")
	form.Set("MASTER_ETCD_INFO", "false")
	form.Set("MASTER_CONTROL_PLANE_RESOURCES", "false")
	form.Set("MASTER_CONTROL_PLANE_RESOURCE_LOGS", "false")
	form.Set("KUBX_CLUSTER_RESOURCES", "false")

	ctx, cancel := context.WithTimeout(ctx, c.cfg.requestTimeout)
	defer cancel()
	resp, err := c.do(ctx, http.MethodPost, jenkinsJobURL+"/buildWithParameters", form)
	if err != nil {
		return "", fmt.Errorf("trigger failed: %w", err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusFound:
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 300))
		return "", fmt.Errorf("trigger failed: status=%d body=%q", resp.StatusCode, body)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", errors.New("trigger succeeded but Jenkins did not return a queue Location header")
	}
	return strings.TrimRight(location, "/") + "/", nil
}

type queueExecutable struct {
	URL    string `json:"url"`
	Number int    `json:"number"`
}

type queueItem struct {
	Cancelled  bool             `json:"cancelled"`
	Executable *queueExecutable `json:"executable"`
}

func (c *jenkinsClient) fetchQueueItem(ctx context.Context, queueURL string) (*queueItem, error) {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.requestTimeout)
	defer cancel()
	resp, err := c.do(ctx, http.MethodGet, strings.TrimRight(queueURL, "/")+"/api/json", nil)
	if err != nil {
		return nil, fmt.Errorf("queue API request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("queue API request failed: status=%d", resp.StatusCode)
	}
	var item queueItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, fmt.Errorf("unexpected Jenkins queue API response shape: %w", err)
	}
	return &item, nil
}

func (c *jenkinsClient) waitForBuildAssignment(ctx context.Context, queueURL string) (string, int, error) {
	for {
		item, err := c.fetchQueueItem(ctx, queueURL)
		if err != nil {
			return "", 0, err
		}
		if item.Cancelled {
			return "", 0, errors.New("the Jenkins queue item was cancelled before a build number was assigned")
		}
		if item.Executable != nil && item.Executable.URL != "" && item.Executable.Number != 0 {
			return strings.TrimRight(item.Executable.URL, "/") + "/", item.Executable.Number, nil
		}
		select {
		case <-ctx.Done():
			return "", 0, ctx.Err()
		case <-time.After(c.cfg.queueInterval):
		}
	}
}

type buildInfo struct {
	Result *string `json:"result"`
}

func (c *jenkinsClient) fetchBuildInfo(ctx context.Context, buildURL string) (*buildInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.requestTimeout)
	defer cancel()
	tree := "number,url,result,building,duration,estimatedDuration,timestamp"
	resp, err := c.do(ctx, http.MethodGet, strings.TrimRight(buildURL, "/")+"/api/json?tree="+tree, nil)
	if err != nil {
		return nil, fmt.Errorf("build API request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("build API request failed: status=%d", resp.StatusCode)
	}
	var build buildInfo
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return nil, fmt.Errorf("unexpected Jenkins build API response shape: %w", err)
	}
	return &build, nil
}

func (c *jenkinsClient) waitForBuildCompletion(ctx context.Context, buildURL string) (*buildInfo, error) {
	for {
		build, err := c.fetchBuildInfo(ctx, buildURL)
		if err != nil {
			return nil, err
		}
		if build.Result != nil {
			return build, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(c.cfg.buildInterval):
		}
	}
}

func (c *jenkinsClient) downloadConsoleText(ctx context.Context, buildURL string) (string, error) {
	// Jenkins consoles can be large; give the download more room than a
	// regular API call (same floor as the Python implementation).
	timeout := c.cfg.requestTimeout
	if timeout < 120*time.Second {
		timeout = 120 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	resp, err := c.do(ctx, http.MethodGet, strings.TrimRight(buildURL, "/")+"/consoleText", nil)
	if err != nil {
		return "", fmt.Errorf("consoleText request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("consoleText request failed: status=%d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("consoleText request failed: %w", err)
	}
	return string(body), nil
}

type executeResult struct {
	Results []commandResult `json:"results"`
}

func executeReadonlyCommands(ctx context.Context, cfg serverConfig, clusterIDRaw, commandsRaw any) (*executeResult, error) {
	clusterID, ok := clusterIDRaw.(string)
	if !ok || strings.TrimSpace(clusterID) == "" {
		return nil, errors.New("cluster_id must be a non-empty string")
	}
	commands, err := normalizeCommands(commandsRaw)
	if err != nil {
		return nil, err
	}
	client, err := newJenkinsClient(cfg)
	if err != nil {
		return nil, err
	}
	return client.runCommands(ctx, strings.TrimSpace(clusterID), commands)
}

func (c *jenkinsClient) runCommands(ctx context.Context, clusterID string, commands []string) (*executeResult, error) {
	queueURL, err := c.triggerJob(ctx, clusterID, commands)
	if err != nil {
		return nil, err
	}
	buildURL, buildNumber, err := c.waitForBuildAssignment(ctx, queueURL)
	if err != nil {
		return nil, err
	}
	build, err := c.waitForBuildCompletion(ctx, buildURL)
	if err != nil {
		return nil, err
	}
	consoleText, err := c.downloadConsoleText(ctx, buildURL)
	if err != nil {
		return nil, err
	}

	parsedBlocks := parseConsoleOutput(consoleText)
	results, warnings, err := pairResultsWithRequestedCommands(commands, parsedBlocks)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("the Jenkins build #%d finished with result %s, but no command output blocks were parsed from consoleText", buildNumber, resultString(build))
	}
	if len(warnings) > 0 {
		return nil, errors.New(strings.Join(warnings, "; "))
	}
	return &executeResult{Results: results}, nil
}

func resultString(build *buildInfo) string {
	if build.Result == nil {
		return "null"
	}
	return *build.Result
}
