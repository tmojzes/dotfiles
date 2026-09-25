package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "armada-xo-slack"
	serverVersion = "0.1.0"
)

const toolDescription = "Use this tool to run a command against the armada-xo bot in the " +
	"#armada-xo channel of IBM's enterprise Slack." +
	"Mostly it us used to query satellite location and kubernetes cluster statuses." +
	"The command is posted as a reply under the latest ':thread:debug satellite' thread in that channel (the thread must already " +
	"exist; create it in Slack if the tool reports it missing), the bot's replies in the " +
	"thread are collected and returned. Use the default (production) bot unless you are " +
	"interested in a cluster in stage; for stage clusters pass bot U08DHQA1FG8. The bot " +
	"replies within seconds to about a minute; the tool waits until it goes quiet, so a " +
	"single call is enough.\n\n" + armadaCommandHelp

// slackCommandInput is the tools/call argument shape. The SDK infers the
// input schema from it; jsonschema tags carry the property descriptions.
type slackCommandInput struct {
	Command string `json:"command" jsonschema:"Command to run against the armada-xo bot, without the @bot mention. Example: 'cluster da3kv2hz06cabscghud0'"`
	Bot     string `json:"bot,omitempty" jsonschema:"Optional bot user ID. Defaults to the production bot U08AUQQRPKJ; pass U08DHQA1FG8 for stage clusters"`
}

func slackInputSchema() (*jsonschema.Schema, error) {
	return jsonschema.For[slackCommandInput](nil)
}

func newServer(cfg serverConfig) (*mcp.Server, error) {
	inputSchema, err := slackInputSchema()
	if err != nil {
		return nil, err
	}
	server := mcp.NewServer(&mcp.Implementation{Name: serverName, Version: serverVersion}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "slack_command",
		Description: toolDescription,
		InputSchema: inputSchema,
	}, slackCommandHandler(cfg))
	return server, nil
}

// slackCommandHandler posts the command and collects the bot's replies. Any
// error it returns becomes a tool error result (IsError) per the SDK
// contract — matching the armada server's behavior.
func slackCommandHandler(cfg serverConfig) func(context.Context, *mcp.CallToolRequest, slackCommandInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input slackCommandInput) (*mcp.CallToolResult, any, error) {
		text, err := runSlackCommand(ctx, cfg, input)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: text}},
		}, nil, nil
	}
}

// runSlackCommand validates the input, finds the latest
// ':thread:debug satellite' thread in #armada-xo, posts "<@bot> <command>"
// as a reply in that thread, and returns the bot's newer thread replies
// rendered as text. The server never posts outside #armada-xo and never
// creates channel root messages.
func runSlackCommand(ctx context.Context, cfg serverConfig, input slackCommandInput) (string, error) {
	command := strings.TrimSpace(input.Command)
	if command == "" {
		return "", errors.New("command must be a non-empty string")
	}

	bot := armadaBotProd
	if provided := strings.TrimSpace(input.Bot); provided != "" {
		if provided != armadaBotProd && provided != armadaBotStage {
			return "", fmt.Errorf("bot must be %s (production, default) or %s (stage), got %q", armadaBotProd, armadaBotStage, provided)
		}
		bot = provided
	}

	client, err := newSlackClient(cfg)
	if err != nil {
		return "", err
	}

	threadTS, err := client.latestThreadRoot(ctx, debugThreadTitle)
	if err != nil {
		return "", err
	}
	postTS, err := client.postMessage(ctx, threadTS, "<@"+bot+"> "+command)
	if err != nil {
		return "", err
	}
	replies, timedOut, err := collectReplies(ctx, client, threadTS, postTS, cfg)
	if err != nil {
		return "", err
	}

	text := renderReplies(replies, threadTS)
	if timedOut {
		text = "The reply timeout was reached; these are the replies collected so far, " +
			"the bot may still be working — check the thread manually.\n\n" + text
	}
	return text, nil
}
