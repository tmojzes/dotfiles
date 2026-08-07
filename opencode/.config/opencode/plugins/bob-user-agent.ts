import { readFileSync } from "node:fs";
import type { Plugin } from "@opencode-ai/plugin";

// Bob CLI's gateway client (class Pc in bob.js) sends
//   User-Agent: bobshell/<cli-version>
// on every request to api.us-east.bob.ibm.com. The gateway rejects
// requests with other User-Agents, so spoof it here.
function bobVersion() {
  try {
    const pkg = JSON.parse(
      readFileSync("/usr/local/lib/node_modules/bobshell/package.json", "utf8"),
    );
    return pkg.version;
  } catch {
    return "1.0.6"; // fallback: check with `bob --version`
  }
}

const BOB_USER_AGENT = `bobshell/${bobVersion()}`;

export default (async () => {
  return {
    "chat.headers": async (input, output) => {
      if (input.provider.id !== "bobshell") return;
      output.headers["User-Agent"] = BOB_USER_AGENT;
    },
  };
}) satisfies Plugin;
