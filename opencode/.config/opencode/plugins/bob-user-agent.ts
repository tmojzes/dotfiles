import { execFileSync } from "node:child_process";
import { readFileSync } from "node:fs";
import { Plugin } from "@opencode/plugin";

// Bob CLI's gateway client (class Pc in bob.js) sends
//   User-Agent: bobshell/<cli-version>
// on every request to api.us-east.bob.ibm.com. The gateway rejects
// requests with other User-Agents, so spoof it here.
function bobVersion() {
  // Authoritative: follows the same PATH the gateway traffic comes from.
  // Output is `<version>\ncommit: <hash>` — take the first line only, and
  // only accept it when it looks like a version.
  try {
    const first = execFileSync("bob", ["--version"], {
      encoding: "utf8",
    }).split("\n")[0].trim();
    if (/^\d/.test(first)) return first;
  } catch {
    // fall through to known install locations
  }
  const roots = [
    `${process.env.HOME}/.npm-global/lib/node_modules`,
    "/usr/local/lib/node_modules",
  ];
  for (const root of roots) {
    try {
      const pkg = JSON.parse(
        readFileSync(`${root}/bobshell/package.json`, "utf8"),
      );
      return pkg.version;
    } catch {
      // try next candidate
    }
  }
  return "1.0.6"; // last resort: check with `bob --version`
}

const BOB_USER_AGENT = `bobshell/${bobVersion()}`;

export default Plugin.define({
  id: "bob-user-agent",
  async setup(ctx) {
    await ctx.session.hook(
      "model.request",
      (event) => {
        event.headers["User-Agent"] = BOB_USER_AGENT;
      },
      { providerID: "bobshell" },
    );
  },
});
