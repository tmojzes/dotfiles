---
name: CI Failure
description: Investigate why a CI check run failed on github.ibm.com by fetching task results and fetching Loki logs for failed tasks
---

Use this skill when the user shares a GitHub CI check URL or check run ID and wants to understand why it failed.

## Tools

- `gh` CLI authenticated against `github.ibm.com` (use `--hostname github.ibm.com`)
- `logcli` for Loki log queries; bearer token is in `~/keys/loki-bearer-token`

## Workflow

### 1. Parse the input

Accept any of:
- A full checks URL: `https://github.ibm.com/<org>/<repo>/pull/<N>/checks?check_run_id=<ID>`
- Just a check run ID + repo context from the user

Extract `org`, `repo`, and `check_run_id`.

### 2. Fetch the check run

```bash
gh api --hostname github.ibm.com \
  /repos/<org>/<repo>/check-runs/<check_run_id>
```

The response body is a Markdown document embedded in the `output.text` field.
Parse it for task entries. Each entry looks like:

```markdown
## **<task-name>** - ❌ Failed <duration>

**Task UID:** <uuid>
**Succeeded:** False
**Reason:** Failed
**Message:** "<step-name>" exited with code <N>: <short message>
**Last transition:** <RFC1123 timestamp>

<LogCLI command block>
```

Extract:
- All **failed** tasks (`:x:` / `❌` in the heading, or `**Succeeded:** False`)
- Their `Task UID`, `Last transition` time, and the embedded `logcli` command
- The pipeline-wide `--from` and `--to` timestamps (present in every logcli command)

### 3. Fetch logs for each failed task

Read the bearer token:
```bash
LOKI_TOKEN=$(cat ~/keys/loki-bearer-token)
```

Run the logcli command extracted from the check run body, adding the bearer token:

```bash
logcli query \
  '{tekton_dev_taskRunUID="<uuid>"} | json | line_format "{{.log}}"' \
  --from=<from> \
  --to=<to> \
  --addr=https://loki.iks-tekton-ci.us-east.containers.appdomain.cloud \
  --bearer-token="$LOKI_TOKEN" \
  --quiet \
  --output=raw \
  --forward \
  --limit=0
```

If there are multiple failed tasks, fetch them concurrently (run in background with `&`, then `wait`) or sequentially — whichever is cleaner for the shell tool.

### 4. Analyse and report

For each failed task, report:
1. **Task name** and **exit message** (from the check run body)
2. **Root cause** — the first error / fatal line in the logs, or the last non-trivial output before the step exited
3. **Relevant log excerpt** (keep it concise; trim repetitive output)

Conclude with a **Summary** section stating what failed and why, and suggest a fix if the cause is clear.

## Notes

- The Loki address is always `https://loki.iks-tekton-ci.us-east.containers.appdomain.cloud` for this CI system.
- `--limit=0` fetches all log lines; do not change this.
- `--forward` preserves chronological order.
- If logcli returns no lines, the task UID may be wrong or the retention window expired — report this clearly.
- Do not print the raw bearer token in output.
