// Command armada-mcp is an MCP server that runs readonly kubectl/oc commands
// via a Jenkins job, built on the official MCP Go SDK.
//
// Ported from the former Python implementation
// (opencode/.config/opencode/mcp/armada.py).
//
// Environment contract (registered in the global opencode.json under
// mcp.servers.armada.environment; hand-exported from ~/keys):
//
//	JENKINS_USER     Jenkins user for the armada job
//	JENKINS_API_KEY  Jenkins API token (~/keys/jenkins-api-key)
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// serverConfig holds the polling and HTTP timing knobs (CLI flags).
type serverConfig struct {
	queueInterval  time.Duration
	buildInterval  time.Duration
	requestTimeout time.Duration
}

func main() {
	var cfg serverConfig
	flag.DurationVar(&cfg.queueInterval, "queue-interval", 5*time.Second, "time between Jenkins queue polls")
	flag.DurationVar(&cfg.buildInterval, "build-interval", 10*time.Second, "time between Jenkins build polls")
	flag.DurationVar(&cfg.requestTimeout, "request-timeout", 30*time.Second, "HTTP timeout for Jenkins requests")
	flag.Parse()

	server, err := newServer(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "armada-mcp:", err)
		os.Exit(1)
	}
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintln(os.Stderr, "armada-mcp:", err)
		os.Exit(1)
	}
}
