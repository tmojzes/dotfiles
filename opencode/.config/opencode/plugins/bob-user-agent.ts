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
    // every body representation seen in the wild: plain string, Buffer/
    // Uint8Array, fetch Request (clone + text), or a text()-bearing object.
    await ctx.session.hook(
      "http.request",
      async (event) => {
        const req = (event as { request?: unknown }).request as
          | (Record<string, unknown> & { clone?: unknown })
          | undefined;
        if (!req) return;

        let bodyText: string | undefined;
        let bodyKind = "unknown";
        const raw = req.body;
        if (typeof raw === "string") {
          bodyText = raw;
          bodyKind = "string";
        } else if (raw instanceof Uint8Array) {
          bodyText = new TextDecoder().decode(raw);
          bodyKind = "uint8array";
        } else if (typeof req.clone === "function") {
          try {
            bodyText = await (req as Request).clone().text();
            bodyKind = "fetch-clone";
          } catch {
            bodyKind = "clone-failed";
          }
        } else if (typeof req.text === "function") {
          try {
            bodyText = await (req.text() as Promise<string>);
            bodyKind = "text()";
          } catch {
            bodyKind = "text-failed";
          }
        }
        if (debug) {
          const proto = Object.getPrototypeOf(event) as object | null;
          const reqProto = req ? Object.getPrototypeOf(req) : null;
          dump("http.request.shape", {
            eventProto: proto ? Object.getOwnPropertyNames(proto) : null,
            requestKind: req?.constructor?.name,
            requestProto: reqProto ? Object.getOwnPropertyNames(reqProto) : null,
            requestKeys: req ? Object.keys(req) : null,
            bodyKind,
            bodyType: typeof raw,
          });
        }
        if (!bodyText || !bodyText.includes('"strict"')) return;

        let parsed: {
          tools?: Array<{ function?: Record<string, unknown> }>;
          messages?: Array<Record<string, unknown>>;
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
        // The gateway also rejects reasoning parts replayed on assistant
        // messages (422 "unexpected property" at body.messages[*].reasoning).
        for (const message of parsed.messages ?? []) {
          if (!message || typeof message !== "object") continue;
          for (const key of ["reasoning", "reasoning_content"]) {
            if (key in message) {
              delete message[key];
              changed = true;
            }
          }
        }
        if (!changed) return;

        const newBody = JSON.stringify(parsed);
        if (bodyKind === "string") {
          req.body = newBody;
        } else if (raw instanceof Uint8Array) {
          req.body = new TextEncoder().encode(newBody);
        } else {
          event.request = new Request(req.url as string, {
            method: req.method as string,
            headers: req.headers as HeadersInit,
            body: newBody,
          });
        }
        dump("strict-stripped", { bodyKind, tools: parsed.tools?.length });
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
      async (event) => {
        const resp = (event as { response?: Response }).response;
        let body: string | undefined;
        try {
          if (typeof resp?.clone === "function") {
            body = (await resp.clone().text()).slice(0, 1500);
          }
        } catch {
          // streaming bodies may not clone; ignore
        }
        dump("http.response", { status: resp?.status, body });
      },
      { providerID: "bobshell" },
    );
  },
};
