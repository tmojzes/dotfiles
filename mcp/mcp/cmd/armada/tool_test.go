package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func testConfig() serverConfig {
	return serverConfig{queueInterval: 10 * time.Millisecond, buildInterval: 10 * time.Millisecond, requestTimeout: 5 * time.Second}
}

// connectTestSession wires a real SDK client to the server over in-memory
// transports, covering the initialize handshake and tools/list.
func connectTestSession(t *testing.T, cfg serverConfig) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()

	server, err := newServer(cfg)
	if err != nil {
		t.Fatalf("newServer: %v", err)
	}
	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0.0.0"}, nil)
	session, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { session.Close() })
	return session
}

func TestToolAdvertised(t *testing.T) {
	session := connectTestSession(t, testConfig())

	found := false
	for tool, err := range session.Tools(context.Background(), nil) {
		if err != nil {
			t.Fatalf("tools iteration: %v", err)
		}
		if tool.Name == "armada_kubectl" {
			found = true
		}
	}
	if !found {
		t.Fatal("tools/list does not advertise armada_kubectl")
	}
}

func TestToolCallMissingCredentials(t *testing.T) {
	t.Setenv("JENKINS_USER", "")
	t.Setenv("JENKINS_API_KEY", "")

	session := connectTestSession(t, testConfig())
	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "armada_kubectl",
		Arguments: map[string]any{"cluster_id": "prod-syd01", "commands": []any{"get nodes"}},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if !res.IsError {
		t.Fatalf("want a tool error for missing credentials, got %+v", res)
	}
	text := res.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(text, "JENKINS_USER") {
		t.Fatalf("tool error text = %q, want it to mention JENKINS_USER", text)
	}
}

const fakeConsole = `Started by user jenkins
Found kubeconfig in prod-syd01 secret
11:22:33 CUSTOM KUBX KUBECTL COMMAND: get nodes
CUSTOM KUBX KUBECTL:
11:22:33 NAME    STATUS
11:22:34 node1   Ready

--------------------
11:22:35 CUSTOM KUBX KUBECTL COMMAND: get pods -A
11:22:36 NS    POD

====================
Finished: SUCCESS
`

// TestToolCallEndToEnd drives the full flow against a fake Jenkins: trigger
// (asserting auth and parameters), queue assignment, build completion, and
// consoleText parsing.
func TestToolCallEndToEnd(t *testing.T) {
	t.Setenv("JENKINS_USER", "testuser")
	t.Setenv("JENKINS_API_KEY", "testtoken")

	jobPath := jenkinsJobPath
	var fake *httptest.Server
	fake = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}
		switch {
		case r.Method == http.MethodPost && r.URL.Path == jobPath+"/buildWithParameters":
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if got := r.PostForm.Get("CUSTOM_KUBX_KUBECTL"); got != "get nodes; get pods -A" {
				http.Error(w, fmt.Sprintf("unexpected commands %q", got), http.StatusBadRequest)
				return
			}
			if got := r.PostForm.Get("CLUSTER"); got != "prod-syd01" {
				http.Error(w, fmt.Sprintf("unexpected cluster %q", got), http.StatusBadRequest)
				return
			}
			w.Header().Set("Location", fake.URL+"/queue/item/1/")
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodGet && r.URL.Path == "/queue/item/1/api/json":
			fmt.Fprint(w, `{"executable":{"url":"`+fake.URL+jobPath+`/1/","number":1}}`)
		case r.Method == http.MethodGet && r.URL.Path == jobPath+"/1/api/json":
			fmt.Fprint(w, `{"number":1,"result":"SUCCESS","building":false}`)
		case r.Method == http.MethodGet && r.URL.Path == jobPath+"/1/consoleText":
			fmt.Fprint(w, fakeConsole)
		default:
			http.Error(w, "unexpected "+r.Method+" "+r.URL.Path, http.StatusNotFound)
		}
	}))
	t.Cleanup(fake.Close)

	originalJobURL := jenkinsJobURL
	jenkinsJobURL = fake.URL + jobPath
	t.Cleanup(func() { jenkinsJobURL = originalJobURL })

	session := connectTestSession(t, testConfig())
	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "armada_kubectl",
		Arguments: map[string]any{"cluster_id": "prod-syd01", "commands": []any{"kubectl get nodes", "get pods -A"}},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool call failed: %+v", res)
	}
	text := res.Content[0].(*mcp.TextContent).Text
	for _, want := range []string{"get nodes", "NAME    STATUS", "node1   Ready", "get pods -A", "NS    POD"} {
		if !strings.Contains(text, want) {
			t.Errorf("result text missing %q:\n%s", want, text)
		}
	}
}
