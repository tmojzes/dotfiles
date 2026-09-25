// Command armada-mcp is an MCP server that runs readonly kubectl/oc commands
// via a Jenkins job.
//
// Status: stub. Port the implementation from
// opencode/.config/opencode/mcp/armada.py (kept until this port lands).
//
// Environment contract (registered in the global opencode.json under
// mcp.servers.armada.environment; hand-exported from ~/keys):
//
//	JENKINS_USER     Jenkins user for the armada job
//	JENKINS_API_KEY  Jenkins API token (~/keys/jenkins-api-key)
package main

import (
	"fmt"
	"os"
)

func main() {
	if os.Getenv("JENKINS_USER") == "" || os.Getenv("JENKINS_API_KEY") == "" {
		fmt.Fprintln(os.Stderr, "armada-mcp: JENKINS_USER and JENKINS_API_KEY must be set")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "armada-mcp: not implemented yet (stub)")
	os.Exit(1)
}
