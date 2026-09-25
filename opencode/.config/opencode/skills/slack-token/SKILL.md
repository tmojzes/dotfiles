---
name: Slack Token
description: Retrieve working Slack session credentials (xoxc- token + d cookie) from Chrome storage or the Slack app browser session
metadata:
  opencode/autoinvoke: false
---

Retrieve the user's Slack session credentials. Slack API calls need **both** an
`xoxc-` client session token (sent as `Authorization: Bearer`) **and** the `d`
cookie value (sent as `Cookie: d=...`). Either one alone returns
`invalid_auth`. There are two methods, tried in order:

1. **Chrome on-disk extraction** (preferred — no browser connection required; covers both normal Slack and enterprise grids like `ibm.enterprise.slack.com`)
2. **Embedded browser / localStorage** (fallback — requires the desktop app browser feature)

---

## Method 1: Chrome on-disk extraction (no browser required)

Reads the `d` cookie from Chrome's SQLite cookie store (decrypting with the
macOS Keychain key) and the `xoxc-` token from Chrome's localStorage leveldb.
Check **all profiles** — each Chrome profile holds a different Slack workspace
session (e.g. `Default` = personal workspace, `Profile 1` = IBM enterprise
grid).

### Step 1 — Extract and decrypt the `d` cookie (all profiles)

Run this Python script via the shell tool:

```python
import sqlite3, shutil, os, subprocess, hashlib, re, glob

result = subprocess.run(
    ["security", "find-generic-password", "-w", "-s", "Chrome Safe Storage"],
    capture_output=True, text=True
)
password = result.stdout.strip().encode()
key = hashlib.pbkdf2_hmac('sha1', password, b'saltysalt', 1003, dklen=16)
iv_hex = '20' * 16  # AES-128-CBC, IV = 16 space chars

for db in glob.glob(os.path.expanduser("~/Library/Application Support/Google/Chrome/*/Cookies")):
    profile = db.split("/")[-2]
    shutil.copy2(db, "/tmp/slack_cookies_tmp.db")
    conn = sqlite3.connect("/tmp/slack_cookies_tmp.db")
    rows = conn.execute(
        "SELECT host_key, name, encrypted_value FROM cookies WHERE host_key LIKE '%slack.com' AND name='d'"
    ).fetchall()
    conn.close()
    for host, name, enc in rows:
        enc_bytes = bytes(enc)
        if enc_bytes[:3] != b'v10':
            continue
        with open('/tmp/slack_c_enc.bin', 'wb') as f:
            f.write(enc_bytes[3:])
        r = subprocess.run(
            ["openssl", "enc", "-aes-128-cbc", "-d", "-nosalt",
             "-K", key.hex(), "-iv", iv_hex, "-nopad", "-in", "/tmp/slack_c_enc.bin"],
            capture_output=True
        )
        if not r.stdout:
            continue
        pad = r.stdout[-1]
        value = r.stdout[:-pad].decode('utf-8', errors='replace')
        tokens = re.findall(r'xox[a-z]-[0-9A-Za-z%-]+', value)
        if tokens:
            print(f"{profile}|{host}|{tokens[0]}")
```

### Step 2 — Extract the `xoxc-` token from localStorage leveldb (all profiles)

```bash
for prof in "Default" "Profile 1" "Profile 2" "Profile 3"; do
  d=~/Library/Application\ Support/Google/Chrome/$prof/Local\ Storage/leveldb
  [ -d "$d" ] || continue
  hits=$(strings "$d"/*.ldb "$d"/*.log 2>/dev/null | grep -o 'xoxc-[0-9A-Za-z][0-9A-Za-z-]*' | sort -u)
  if [ -n "$hits" ]; then echo "== $prof =="; echo "$hits"; fi
done
```

Leveldb keeps old values, so a profile can yield **multiple** `xoxc-` tokens
(including stale ones from expired sessions). Test each candidate in Step 3.

### Step 3 — Verify with `auth.test` and identify the workspace

```bash
curl -sS -H "Authorization: Bearer $XOXC" -H "Cookie: d=$XOXD" \
  "https://<host>/api/auth.test"
```

- `<host>` is `api.slack.com` for normal workspaces, or the enterprise grid's
  vanity domain (e.g. `ibm.enterprise.slack.com` for IBM) — both work for
  enterprise grids.
- `ok: true` returns `team`, `team_id`, `user`, `enterprise_id` — use this to
  pick the profile whose tokens match the workspace the user asked about.
- Use **matching** `xoxc` + `xoxd` from the **same profile**.
- The `d` cookie value works URL-encoded or URL-decoded (both verified).

### Step 4 — Save to `~/keys/` and export

Once the correct pair is identified, save them to the canonical key files:

```bash
echo -n "<xoxc-token>" > ~/keys/slack-xoxc
echo -n "<xoxd-value>" > ~/keys/slack-d-cookie
```

Then provide the export commands to the user so they can reload the MCP server environment:

```bash
export SLACK_XOXC=$(cat ~/keys/slack-xoxc)
export SLACK_XOXD=$(cat ~/keys/slack-d-cookie)
```

**Important:** The `slack` MCP server reads `SLACK_XOXC`/`SLACK_XOXD` from the environment it was launched with (`{env:SLACK_XOXC}` in `opencode.json`). If OpenCode was started before these were exported, the MCP server process won't see the new values until OpenCode is restarted. Tell the user to:
1. Export the vars in the shell that launched OpenCode
2. Restart OpenCode (or just the MCP service if hot-reload is available)

### Failure modes

| Symptom | Fix |
|---|---|
| `security` command returns nothing | Chrome Safe Storage entry missing — try Method 2 |
| Decryption produces garbage (no `xox` token) | Wrong profile — try other profiles under `~/Library/Application Support/Google/Chrome/` |
| No `.slack.com` rows in the cookie DB | User isn't logged into Slack in Chrome — use Method 2 or ask them to log in in Chrome |
| `xoxc` found but `auth.test` says `invalid_auth` | Stale/expired session — try the next candidate token, re-login in Chrome, or use Method 2 |
| No `xoxc` in leveldb, only `xoxd` | `xoxd` alone is NOT enough for API calls — the user must (re)login so `localConfig_v2` is written, or use Method 2 |
| `invalid_auth` with both tokens | Tokens are from different profiles/workspaces, or the session expired — re-extract and verify |

---

## Method 2: Embedded browser / localStorage (fallback)

Use this if Chrome storage isn't available or the user is logged into Slack
only in the OpenCode desktop app's embedded browser.

### Prerequisites

- OpenCode desktop app installed with the **browser experimental feature** enabled.
- User must be logged into Slack in the embedded browser. The embedded browser is a **separate profile** from Chrome/Safari — existing sessions don't carry over. If not logged in, open `https://app.slack.com/` in the embedded browser and ask them to log in there first.

### Step 1 — Open a Slack tab

Open `https://app.slack.com/` in a new browser tab and wait for it to finish loading. Verify the final URL contains `app.slack.com` (not a sign-in redirect). If the URL contains `/sign_in` or `/ssb/redirect`, the user is not logged in — stop and ask them to log in.

### Step 2 — Extract token via JavaScript

Execute this in the loaded tab:

```js
(() => {
  const out = { tokens: [], dCookie: document.cookie.split(/;\s*/)
    .filter(c => c.startsWith("d=")).map(c => c.slice(2)) };
  // Primary: parse localConfig_v2 (present in most Slack versions)
  try {
    const cfg = JSON.parse(localStorage.getItem("localConfig_v2") || "{}");
    const teams = cfg.teams || {};
    out.tokens.push(...Object.values(teams).map(t => ({
      source: "localConfig_v2",
      name: t.name || t.domain || "(unknown)",
      domain: t.domain || "",
      token: t.token || null,
    })));
  } catch (_) {}
  // Fallback: scan all localStorage keys for xoxc- tokens
  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i);
    try {
      const val = localStorage.getItem(key) || "";
      const matches = val.match(/xoxc-[0-9A-Za-z\-]+/g);
      if (matches) out.tokens.push(...matches.map(t => ({ source: "scan", key, token: t })));
    } catch (_) {}
  }
  return out;
})()
```

### Step 3 — Report the result

- Present each `xoxc-` token labelled by workspace name/domain, **together with the `d` cookie value** — the API needs both (see Method 1, Step 3 for the request shape).
- If multiple workspaces are returned, list all of them.
- If no tokens are returned, the embedded browser is not logged in — fall back to Method 1 or ask the user to log in.

---

## API usage (verified)

```bash
curl -H "Authorization: Bearer $XOXC" -H "Cookie: d=$XOXD" \
  "https://<host>/api/<method>"
```

- Works on `api.slack.com` and on enterprise vanity domains (verified with
  `ibm.enterprise.slack.com`).
- `xoxd-` alone — as Bearer, form `token=` param, encoded or decoded — returns
  `invalid_auth`. Don't bother trying it.

## Token types

| Token prefix | Type | Source | Usable for API calls |
|---|---|---|---|
| `xoxc-` | Client session token | localStorage (`localConfig_v2`) | Yes — **with** `d` cookie as well |
| `xoxd-` | Cookie-based session token | Chrome `d` cookie | Only as the `Cookie: d=` companion to `xoxc-` |

Both expire when the Slack session ends; re-run this skill to refresh them.

## Security note

These tokens grant full access to the Slack workspace as the signed-in user. Do not log, store, or transmit them beyond returning them to the user in this session.
