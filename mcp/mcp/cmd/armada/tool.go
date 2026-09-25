package main

import (
	"context"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "jenkins-readonly-kubectl"
	serverVersion = "0.2.0"
)

const toolDescription = "Use this tool for readonly investigation of a specific armada cluster. " +
	"It runs kubectl/oc get, describe, logs, top, api-resources, or version against the " +
	"cluster_id you provide. Prefer this tool when the user asks for fresh, realtime " +
	"cluster inspection, or explicitly says not to use the kubernetes MCP server. " +
	"To run multiple commands, pass them as multiple strings in commands, or as one " +
	"string separated by ';' or newlines. Shell pipes are NOT supported. As this runs " +
	"on the backend as a Jenkins job it can take a long time to finish, about 5 " +
	"minutes; instead of using the tool multiple times you should run it once with " +
	"all the commands you need."

// armadaKubectlInput is the tools/call argument shape. The SDK infers the
// input schema from it (both fields required); jsonschema tags carry the
// property descriptions.
type armadaKubectlInput struct {
	ClusterID string   `json:"cluster_id" jsonschema:"The opaque Armada cluster ID to pass as the CLUSTER Jenkins parameter. This is NOT the human-readable carrier name (e.g. not 'prod-lon04-carrier106'). It is a short alphanumeric string of 16-20 lowercase letters and digits, for example 'a1b2c3d4e5f6g7h8'. Obtain it from cluster listings or the user — do not pass a carrier hostname."`
	Commands  []string `json:"commands" jsonschema:"List of readonly kubectl subcommands. Each entry may optionally start with 'kubectl' or 'oc', but only get, describe, logs, top, api-resources, and version are allowed. Multiple commands may be provided either as separate array items or combined in a single string separated by ';' or newlines. Example: ['kubectl get nodes', 'kubectl get pods -A']"`
}

// armadaInputSchema infers the schema from the input struct and carries over
// the minItems constraint the former Python server advertised.
func armadaInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[armadaKubectlInput](nil)
	if err != nil {
		return nil, err
	}
	schema.Properties["commands"].MinItems = jsonschema.Ptr(1)
	return schema, nil
}

func newServer(cfg serverConfig) (*mcp.Server, error) {
	inputSchema, err := armadaInputSchema()
	if err != nil {
		return nil, err
	}
	server := mcp.NewServer(&mcp.Implementation{Name: serverName, Version: serverVersion}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "armada_kubectl",
		Description: toolDescription,
		InputSchema: inputSchema,
	}, armadaKubectlHandler(cfg))
	return server, nil
}

// armadaKubectlHandler runs the readonly commands via the Jenkins job. Any
// error it returns becomes a tool error result (IsError) per the SDK
// contract — matching the former Python behavior.
func armadaKubectlHandler(cfg serverConfig) func(context.Context, *mcp.CallToolRequest, armadaKubectlInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input armadaKubectlInput) (*mcp.CallToolResult, any, error) {
		payload, err := executeReadonlyCommands(ctx, cfg, input.ClusterID, input.Commands)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: renderToolText(payload)}},
		}, nil, nil
	}
}

// renderToolText formats results the way the former Python server did: a
// single command returns its bare output, multiple commands are prefixed
// with their command line and joined by a blank line.
func renderToolText(payload *executeResult) string {
	results := payload.Results
	if len(results) == 1 {
		return results[0].Output
	}
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
	return strings.Join(chunks, "\n\n")
}
