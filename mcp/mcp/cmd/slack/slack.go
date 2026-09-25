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
	"strconv"
	"strings"
	"time"
)

const (
	// slackAPIBase is IBM's enterprise Slack grid; the standard Slack Web
	// API is served from the grid's vanity domain.
	slackAPIBase = "https://ibm.enterprise.slack.com/api"

	// armadaChannelID is #armada-xo, the private channel where the
	// armada-xo bots are invoked.
	armadaChannelID = "C53AJ95TP"

	// The two armada-xo bot user IDs seen in the channel: the production
	// bot and the stage bot.
	armadaBotProd  = "U08AUQQRPKJ"
	armadaBotStage = "U08DHQA1FG8"

	// debugThreadTitle is the title of the thread in #armada-xo that bot
	// commands are posted under. Several threads share the title over time;
	// the server always picks the latest one.
	debugThreadTitle = ":thread:debug satellite"

	// replyPageLimit is the maximum messages per conversations.replies call;
	// historyPageLimit likewise for conversations.history, which is paged
	// only up to maxHistoryPages pages when searching for the thread root.
	replyPageLimit   = 200
	historyPageLimit = 200
	maxHistoryPages  = 5

	maxRateLimitRetries    = 2
	maxRateLimitRetryAfter = 10 * time.Second
)

// slackMessage is the subset of a Slack message the server cares about: its
// timestamp (unique per channel, and the thread reply ordering key) and text.
type slackMessage struct {
	Ts   string `json:"ts"`
	Text string `json:"text"`
}

// slackResponse is the envelope every Slack Web API method wraps its payload
// in; ok=false carries a machine-readable error code.
type slackResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type slackClient struct {
	httpClient *http.Client
	xoxc       string
	xoxd       string
	cfg        serverConfig
}

// newSlackClient reads the hand-exported Slack session credentials from the
// environment. Missing credentials are a per-call tool error, not a startup
// failure, so the server stays connectable without them.
func newSlackClient(cfg serverConfig) (*slackClient, error) {
	xoxc := os.Getenv("SLACK_XOXC")
	xoxd := os.Getenv("SLACK_XOXD")
	if xoxc == "" || xoxd == "" {
		return nil, errors.New("SLACK_XOXC and SLACK_XOXD must be set (session credentials from the slack-token skill, exported from ~/keys)")
	}
	return &slackClient{
		httpClient: &http.Client{Timeout: cfg.requestTimeout},
		xoxc:       xoxc,
		xoxd:       xoxd,
		cfg:        cfg,
	}, nil
}

// postMessage posts text as a reply in the given thread of #armada-xo and
// returns the timestamp of the created message. The thread_ts parameter
// makes the message a thread reply; the server never posts channel roots.
func (c *slackClient) postMessage(ctx context.Context, threadTS, text string) (string, error) {
	form := url.Values{}
	form.Set("channel", armadaChannelID)
	form.Set("thread_ts", threadTS)
	form.Set("text", text)

	payload, err := c.do(ctx, "chat.postMessage", http.MethodPost, form)
	if err != nil {
		return "", err
	}
	var result struct {
		slackResponse
		Ts string `json:"ts"`
	}
	if err := json.Unmarshal(payload, &result); err != nil {
		return "", fmt.Errorf("unexpected chat.postMessage response shape: %w", err)
	}
	if !result.OK {
		return "", slackAPIError(result.Error)
	}
	if result.Ts == "" {
		return "", errors.New("chat.postMessage succeeded but returned no message timestamp")
	}
	return result.Ts, nil
}

// fetchReplies returns all thread replies for the given root message,
// following pagination cursors. The Slack API includes the root message
// itself; callers filter it out.
func (c *slackClient) fetchReplies(ctx context.Context, rootTS string) ([]slackMessage, error) {
	var messages []slackMessage
	cursor := ""
	for {
		params := url.Values{}
		params.Set("channel", armadaChannelID)
		params.Set("ts", rootTS)
		params.Set("limit", strconv.Itoa(replyPageLimit))
		if cursor != "" {
			params.Set("cursor", cursor)
		}

		payload, err := c.do(ctx, "conversations.replies", http.MethodGet, params)
		if err != nil {
			return nil, err
		}
		var result struct {
			slackResponse
			Messages []slackMessage `json:"messages"`
			Metadata *struct {
				NextCursor string `json:"next_cursor"`
			} `json:"response_metadata"`
		}
		if err := json.Unmarshal(payload, &result); err != nil {
			return nil, fmt.Errorf("unexpected conversations.replies response shape: %w", err)
		}
		if !result.OK {
			return nil, slackAPIError(result.Error)
		}
		messages = append(messages, result.Messages...)
		if result.Metadata == nil || result.Metadata.NextCursor == "" {
			return messages, nil
		}
		cursor = result.Metadata.NextCursor
	}
}

// fetchChannelHistory returns the most recent #armada-xo messages, newest
// first, following pagination cursors up to maxHistoryPages pages — enough
// to find the latest debug thread without walking the whole channel.
func (c *slackClient) fetchChannelHistory(ctx context.Context) ([]slackMessage, error) {
	var messages []slackMessage
	cursor := ""
	for page := 0; page < maxHistoryPages; page++ {
		params := url.Values{}
		params.Set("channel", armadaChannelID)
		params.Set("limit", strconv.Itoa(historyPageLimit))
		if cursor != "" {
			params.Set("cursor", cursor)
		}

		payload, err := c.do(ctx, "conversations.history", http.MethodGet, params)
		if err != nil {
			return nil, err
		}
		var result struct {
			slackResponse
			Messages []slackMessage `json:"messages"`
			Metadata *struct {
				NextCursor string `json:"next_cursor"`
			} `json:"response_metadata"`
		}
		if err := json.Unmarshal(payload, &result); err != nil {
			return nil, fmt.Errorf("unexpected conversations.history response shape: %w", err)
		}
		if !result.OK {
			return nil, slackAPIError(result.Error)
		}
		messages = append(messages, result.Messages...)
		if result.Metadata == nil || result.Metadata.NextCursor == "" {
			return messages, nil
		}
		cursor = result.Metadata.NextCursor
	}
	return messages, nil
}

// latestThreadRoot returns the ts of the newest #armada-xo message whose
// text matches the given thread title.
func (c *slackClient) latestThreadRoot(ctx context.Context, title string) (string, error) {
	messages, err := c.fetchChannelHistory(ctx)
	if err != nil {
		return "", err
	}
	return findLatestThreadRoot(messages, title)
}

// replySource is the part of slackClient the reply collector depends on; an
// interface keeps the collector unit-testable.
type replySource interface {
	fetchReplies(ctx context.Context, rootTS string) ([]slackMessage, error)
}

// do invokes a Slack Web API method. GET sends the params as the query
// string, POST as a form body. Session credentials are attached the same way
// the Slack web client sends them: Bearer xoxc plus the d cookie — either
// one alone gets invalid_auth.
func (c *slackClient) do(ctx context.Context, method, httpMethod string, params url.Values) ([]byte, error) {
	rawURL := slackAPIBase + "/" + method
	for attempt := 0; ; attempt++ {
		req, err := c.newRequest(ctx, rawURL, httpMethod, params)
		if err != nil {
			return nil, err
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("slack API request failed: %w", err)
		}

		// Honor a rate-limit response instead of failing the call; the
		// reply poller fires a request every few seconds.
		if resp.StatusCode == http.StatusTooManyRequests && attempt < maxRateLimitRetries {
			retryAfter := c.cfg.pollInterval
			if v, convErr := strconv.Atoi(resp.Header.Get("Retry-After")); convErr == nil && v > 0 {
				retryAfter = time.Duration(v) * time.Second
			}
			if retryAfter > maxRateLimitRetryAfter {
				retryAfter = maxRateLimitRetryAfter
			}
			resp.Body.Close()
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryAfter):
			}
			continue
		}

		payload, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("slack API request failed: %w", readErr)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("slack API request failed: status=%d body=%q", resp.StatusCode, truncateForError(payload))
		}
		return payload, nil
	}
}

func (c *slackClient) newRequest(ctx context.Context, rawURL, httpMethod string, params url.Values) (*http.Request, error) {
	if httpMethod == http.MethodGet {
		if params != nil {
			rawURL += "?" + params.Encode()
		}
		req, err := http.NewRequestWithContext(ctx, httpMethod, rawURL, nil)
		if err != nil {
			return nil, err
		}
		c.setAuth(req)
		return req, nil
	}

	var body io.Reader
	if params != nil {
		body = strings.NewReader(params.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, httpMethod, rawURL, body)
	if err != nil {
		return nil, err
	}
	c.setAuth(req)
	if params != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, nil
}

func (c *slackClient) setAuth(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.xoxc)
	req.Header.Set("Cookie", "d="+c.xoxd)
}

// slackAPIError turns a Slack API error code into a descriptive error,
// adding a refresh hint for expired session credentials.
func slackAPIError(code string) error {
	switch code {
	case "invalid_auth", "not_authed", "token_expired", "cookie_rotation_failure":
		return fmt.Errorf("slack API error %q: the session credentials expired; re-run the slack-token skill and re-export SLACK_XOXC/SLACK_XOXD from ~/keys", code)
	case "not_in_channel":
		return fmt.Errorf("slack API error %q: the token's user is not a member of #armada-xo", code)
	case "":
		return errors.New("unknown slack API error (empty error code)")
	default:
		return fmt.Errorf("slack API error: %s", code)
	}
}

// truncateForError caps an error-message body the way the armada server does.
func truncateForError(body []byte) []byte {
	const limit = 300
	if len(body) <= limit {
		return body
	}
	return body[:limit]
}
