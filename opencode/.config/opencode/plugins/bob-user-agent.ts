import { execFileSync } from "node:child_process";

// Bob 2.x gateway quirks (api.us-east.bob.ibm.com, behind Cloudflare):
//
// 1. Only requests whose User-Agent matches bob's own client UA,
//    `bob-shell/<version>`, get through — the old "bobshell/<version>"
//    string from bob 1.x is rejected with a Cloudflare 403 page.
// 2. The request validation rejects the OpenAI `strict` flag on tool
//    definitions: `tools[*].function.strict` is an "unexpected property"
//    (422, application/problem+json). It has to be stripped from every
//    tool before the request goes out.
//
// NOTE: no `@opencode/plugin` import. The SDK's Plugin.define is an identity
// function (dist/promise/plugin.js), so a plain `{ id, setup }` object is a
// valid V2 export — and opencode resolves plugin modules from the plugin
// file's real path (inside this dotfiles repo), where `@opencode/plugin` is
// not installed. Keeping the file import-free avoids that entirely.
function bobVersion() {
  try {
    return execFileSync("bob", ["--version"], { encoding: "utf8" })
      .split("\n")[0]
      .trim();
  } catch {
    return "2.0.5"; // fallback: `bob --version` output
  }
}

export default {
  id: "bob-user-agent",
  async setup(ctx) {
    const userAgent = `bob-shell/${bobVersion()}`;
    const debug = process.env.BOB_PLUGIN_DEBUG === "1";
    let dump: (label: string, data: unknown) => void = () => {};
    if (debug) {
      const { appendFileSync } = await import("node:fs");
      dump = (label, data) => {
        try {
          appendFileSync(
            "/tmp/bob-plugin-debug.log",
            `=== ${label} ${new Date().toISOString()} ===\n${JSON.stringify(data, null, 1).slice(0, 4000)}\n\n`,
          );
        } catch {
          // best effort only
        }
      };
    }
    dump("setup", { userAgent });

    // User-Agent spoofing: applied to every bobshell model request.
    await ctx.session.hook(
      "model.request",
      (event) => {
        event.headers["User-Agent"] = userAgent;
      },
      { providerID: "bobshell" },
    );

    // Strip `strict` from all tools on the serialized request body. Handles
    // both a plain request object (string body) and a fetch Request.
    await ctx.session.hook(
      "http.request",
      async (event) => {
        const req = event.request as
          | (Request & { body?: unknown })
          | undefined;
        if (!req) return;

        let bodyText: string | undefined;
        let plainBody = false;
        if (typeof req.body === "string") {
          bodyText = req.body;
          plainBody = true;
        } else if (typeof req.clone === "function") {
          try {
            bodyText = await req.clone().text();
          } catch {
            return;
          }
        }
        if (!bodyText || !bodyText.includes('"strict"')) return;

        let parsed: {
          tools?: Array<{ function?: Record<string, unknown> }>;
        };
        try {
          parsed = JSON.parse(bodyText);
        } catch {
          return;
        }
        let changed = false;
        for (const tool of parsed.tools ?? []) {
          const fn = tool?.function;
          if (fn && Object.prototype.hasOwnProperty.call(fn, "strict")) {
            delete fn.strict;
            changed = true;
          }
        }
        if (!changed) return;

        const newBody = JSON.stringify(parsed);
        if (plainBody) {
          req.body = newBody;
        } else {
          event.request = new Request(req.url, {
            method: req.method,
            headers: req.headers,
            body: newBody,
          });
        }
        dump("strict-stripped", { tools: parsed.tools?.length });
      },
      { providerID: "bobshell" },
    );

    if (!debug) return;
    // Opt-in wire dumps (BOB_PLUGIN_DEBUG=1 in the service environment).
    const shape = (obj: unknown): unknown => {
      if (!obj || typeof obj !== "object") return obj;
      return Object.fromEntries(
        Object.keys(obj).map((k) => {
          const v = (obj as Record<string, unknown>)[k];
          return [
            k,
            v === null
              ? null
              : typeof v === "object"
                ? `{keys:${Object.keys(v as object).join(",")}}`
                : typeof v,
          ];
        }),
      );
    };
    await ctx.session.hook(
      "http.request",
      (event) => dump("http.request", shape(event)),
      { providerID: "bobshell" },
    );
    await ctx.session.hook(
      "http.response",
      (event) => dump("http.response", shape(event)),
      { providerID: "bobshell" },
    );
  },
};
