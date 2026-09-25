// Command slack-mcp is an MCP server that runs commands against the armada-xo
// bot in the #armada-xo channel of IBM's enterprise Slack grid
// (ibm.enterprise.slack.com), built on the official MCP Go SDK.
//
// The armada-xo bot is addressed by mentioning it in a reply under the
// latest ':thread:debug satellite' thread in #armada-xo; it answers in the
// same thread, usually with several messages in a burst. The server finds
// the thread, posts the command as a reply in it, polls the thread for
// replies newer than its own post, and returns the collected reply text.
// It only ever posts in #armada-xo, never as a channel root message.
//
// Environment contract (registered in the global opencode.json under
// mcp.servers.slack.environment; hand-exported from ~/keys):
//
//	SLACK_XOXC  Slack client session token (xoxc-...), from the slack-token skill
//	SLACK_XOXD  Slack d cookie value (xoxd-...), from the slack-token skill
//
// Both values expire with the Slack session; re-run the slack-token skill and
// re-export them when they do.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// serverConfig holds the reply-polling and HTTP timing knobs (CLI flags).
type serverConfig struct {
	pollInterval   time.Duration
	quietWindow    time.Duration
	replyTimeout   time.Duration
	requestTimeout time.Duration
}

func main() {
	var cfg serverConfig
	flag.DurationVar(&cfg.pollInterval, "poll-interval", 3*time.Second, "time between thread reply polls")
	flag.DurationVar(&cfg.quietWindow, "quiet-window", 10*time.Second, "silent period after the last new reply before the result is returned")
	flag.DurationVar(&cfg.replyTimeout, "reply-timeout", 3*time.Minute, "overall timeout for collecting bot replies")
	flag.DurationVar(&cfg.requestTimeout, "request-timeout", 30*time.Second, "HTTP timeout for Slack API requests")
	flag.Parse()

	server, err := newServer(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "slack-mcp:", err)
		os.Exit(1)
	}
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintln(os.Stderr, "slack-mcp:", err)
		os.Exit(1)
	}
}
