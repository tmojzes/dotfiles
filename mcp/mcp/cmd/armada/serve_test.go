package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func testConfig() serverConfig {
	return serverConfig{queueInterval: time.Second, buildInterval: time.Second, requestTimeout: 5 * time.Second}
}

// TestServeNewlineFraming walks one message of each kind through the serve
// loop: initialize, ping, tools/list, a notification (no response), a
// tools/call that fails on missing Jenkins credentials, and invalid JSON.
func TestServeNewlineFraming(t *testing.T) {
	t.Setenv("JENKINS_USER", "")
	t.Setenv("JENKINS_API_KEY", "")

	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26"}}`,
		`{"jsonrpc":"2.0","id":2,"method":"ping"}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"armada_kubectl","arguments":{"cluster_id":"prod-syd01","commands":["get nodes"]}}}`,
		`not json`,
		``,
	}, "\n")

	var out bytes.Buffer
	if err := serve(context.Background(), strings.NewReader(input), &out, testConfig()); err != nil {
		t.Fatalf("serve: %v", err)
	}

	var responses []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if line == "" {
			continue
		}
		var response map[string]any
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			t.Fatalf("response is not newline-framed JSON: %v (%q)", err, line)
		}
		responses = append(responses, response)
	}
	if len(responses) != 5 {
		t.Fatalf("got %d responses, want 5: %s", len(responses), out.String())
	}

	initResult := responses[0]["result"].(map[string]any)
	if initResult["protocolVersion"] != "2025-03-26" {
		t.Errorf("protocolVersion = %v, want 2025-03-26", initResult["protocolVersion"])
	}
	serverInfo := initResult["serverInfo"].(map[string]any)
	if serverInfo["name"] != serverName {
		t.Errorf("serverInfo.name = %v, want %s", serverInfo["name"], serverName)
	}

	if responses[1]["result"] == nil {
		t.Errorf("ping result missing: %v", responses[1])
	}

	tools := responses[2]["result"].(map[string]any)["tools"].([]any)
	tool := tools[0].(map[string]any)
	if tool["name"] != "armada_kubectl" {
		t.Errorf("tool name = %v, want armada_kubectl", tool["name"])
	}

	call := responses[3]
	if call["error"] != nil {
		t.Fatalf("tools/call returned a protocol error: %v", call)
	}
	callResult := call["result"].(map[string]any)
	if callResult["isError"] != true {
		t.Fatalf("tools/call result isError = %v, want true (missing credentials): %v", callResult["isError"], call)
	}
	callText := callResult["content"].([]any)[0].(map[string]any)["text"].(string)
	if !strings.Contains(callText, "JENKINS_USER") {
		t.Errorf("tool error text = %q, want it to mention JENKINS_USER", callText)
	}

	parseErr := responses[4]["error"].(map[string]any)
	if parseErr["code"].(float64) != -32700 {
		t.Errorf("invalid JSON error code = %v, want -32700", parseErr["code"])
	}
}

func TestServeContentLengthFraming(t *testing.T) {
	body := `{"jsonrpc":"2.0","id":7,"method":"tools/list"}`
	input := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(body), body)

	var out bytes.Buffer
	if err := serve(context.Background(), strings.NewReader(input), &out, testConfig()); err != nil {
		t.Fatalf("serve: %v", err)
	}

	raw := out.String()
	if !strings.HasPrefix(raw, "Content-Length: ") {
		t.Fatalf("response is not content-length framed: %q", raw)
	}
	header, responseBody, found := strings.Cut(raw, "\r\n\r\n")
	if !found {
		t.Fatalf("response has no header/body split: %q", raw)
	}
	declared := strings.TrimSpace(strings.TrimPrefix(header, "Content-Length: "))
	declaredLen, err := strconv.Atoi(declared)
	if err != nil {
		t.Fatalf("Content-Length header is not a number: %q", header)
	}
	if declaredLen != len(responseBody) {
		t.Fatalf("Content-Length = %d, but body is %d bytes", declaredLen, len(responseBody))
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(responseBody), &response); err != nil {
		t.Fatalf("response body is not JSON: %v (%q)", err, responseBody)
	}
	tools := response["result"].(map[string]any)["tools"].([]any)
	if tool := tools[0].(map[string]any); tool["name"] != "armada_kubectl" {
		t.Errorf("tool name = %v, want armada_kubectl", tool["name"])
	}
}
