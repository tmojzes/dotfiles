package main

import (
	"context"
	"fmt"
	"strings"
)

const (
	serverName             = "jenkins-readonly-kubectl"
	serverVersion          = "0.1.0"
	defaultProtocolVersion = "2025-03-26"
)

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type rpcResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

func successResponse(id, result any) *rpcResponse {
	return &rpcResponse{JSONRPC: "2.0", ID: id, Result: result}
}

func errorResponse(id any, rpcErr *rpcError) *rpcResponse {
	return &rpcResponse{JSONRPC: "2.0", ID: id, Error: rpcErr}
}

// handleRequest dispatches one JSON-RPC method. A nil result with a nil
// rpcError means "no response" (notification).
func handleRequest(ctx context.Context, cfg serverConfig, method string, params any) (any, *rpcError) {
	switch method {
	case "initialize":
		protocolVersion := defaultProtocolVersion
		if paramsMap, ok := params.(map[string]any); ok {
			if requested, ok := paramsMap["protocolVersion"].(string); ok && requested != "" {
				protocolVersion = requested
			}
		}
		return map[string]any{
			"protocolVersion": protocolVersion,
			"capabilities":    map[string]any{"tools": map[string]any{"listChanged": false}},
			"serverInfo":      map[string]any{"name": serverName, "version": serverVersion},
		}, nil

	case "notifications/initialized":
		return nil, nil

	case "ping":
		return map[string]any{}, nil

	case "tools/list":
		return map[string]any{"tools": []any{armadaKubectlTool()}}, nil

	case "tools/call":
		return callTool(ctx, cfg, params)
	}
	return nil, &rpcError{Code: -32601, Message: fmt.Sprintf("Method not found: %s", method)}
}

func callTool(ctx context.Context, cfg serverConfig, params any) (any, *rpcError) {
	paramsMap, ok := params.(map[string]any)
	if !ok {
		return nil, &rpcError{Code: -32602, Message: "tools/call params must be an object."}
	}

	name := paramsMap["name"]
	if name != "armada_kubectl" && name != "run_readonly_kubectl_commands" {
		return nil, &rpcError{Code: -32601, Message: fmt.Sprintf("Unknown tool: %#v", name)}
	}

	arguments, ok := paramsMap["arguments"].(map[string]any)
	if !ok {
		return nil, &rpcError{Code: -32602, Message: "tools/call arguments must be an object."}
	}

	payload, err := executeReadonlyCommands(ctx, cfg, arguments["cluster_id"], arguments["commands"])
	if err != nil {
		return buildToolError(err), nil
	}
	return buildToolResult(payload), nil
}

type toolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

func armadaKubectlTool() toolDefinition {
	return toolDefinition{
		Name: "armada_kubectl",
		Description: "Use this tool for readonly investigation of a specific armada cluster. " +
			"It runs kubectl/oc get, describe, logs, top, api-resources, or version against the " +
			"cluster_id you provide. Prefer this tool when the user asks for fresh, realtime " +
			"cluster inspection, or explicitly says not to use the kubernetes MCP server. " +
			"To run multiple commands, pass them as multiple strings in commands, or as one " +
			"string separated by ';' or newlines. Shell pipes are NOT supported. As this runs " +
			"on the backend as a Jenkins job it can take a long time to finish, about 5 " +
			"minutes; instead of using the tool multiple times you should run it once with " +
			"all the commands you need.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"cluster_id": map[string]any{
					"type":        "string",
					"description": "Cluster ID to pass as the CLUSTER Jenkins parameter.",
				},
				"commands": map[string]any{
					"type": "array",
					"description": "List of readonly kubectl subcommands. Each entry may optionally " +
						"start with 'kubectl' or 'oc', but only get, describe, logs, top, " +
						"api-resources, and version are allowed. Multiple commands may be provided " +
						"either as separate array items or combined in a single string separated " +
						"by ';' or newlines. Example: ['kubectl get nodes', 'kubectl get pods -A'].",
					"items":    map[string]any{"type": "string"},
					"minItems": 1,
				},
			},
			"required":             []any{"cluster_id", "commands"},
			"additionalProperties": false,
		},
	}
}

type toolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func buildToolResult(payload *executeResult) map[string]any {
	results := payload.Results
	var text string
	if len(results) == 1 {
		text = results[0].Output
	} else {
		var chunks []string
		for _, result := range results {
			command := strings.TrimSpace(result.Command)
			chunk := result.Output
			if command != "" {
				chunk = strings.TrimSpace(command + "\n" + result.Output)
			}
			if chunk != "" {
				chunks = append(chunks, chunk)
			}
		}
		text = strings.Join(chunks, "\n\n")
	}
	return map[string]any{
		"content": []any{toolContent{Type: "text", Text: text}},
		"isError": false,
	}
}

func buildToolError(err error) map[string]any {
	return map[string]any{
		"content": []any{toolContent{Type: "text", Text: err.Error()}},
		"isError": true,
	}
}
