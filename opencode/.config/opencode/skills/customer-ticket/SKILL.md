---
name: Customer Ticket
description: Investigate an IBM Cloud / Satellite customer support ticket from github.ibm.com/alchemy-containers/customer-tickets. Fetches the ticket via gh CLI, runs armada-xo Slack bot commands for high-level cluster status, and uses armada kubectl for deeper inspection.
---

Use this skill when the user shares a customer ticket URL or issue number from
`github.ibm.com/alchemy-containers/customer-tickets` and wants to understand
or debug the underlying cluster issue.

## Tools

- `gh` CLI authenticated against `github.ibm.com`
- `slack` MCP tool (`tools.slack.slack_command`) — armada-xo bot in #armada-xo
- `armada` MCP tool (`tools.armada.armada_kubectl`) — readonly kubectl on carrier/tugboat clusters
- Jenkins API (via `curl`) for GMI job logs when deeper history is needed

## Workflow

### 1. Fetch the ticket

```bash
# Get the issue body
gh api --hostname github.ibm.com \
  /repos/alchemy-containers/customer-tickets/issues/<N> --jq '.body'

# Get all comments (includes bot analysis, SRE notes, dev notes)
gh api --hostname github.ibm.com \
  /repos/alchemy-containers/customer-tickets/issues/<N>/comments \
  --jq '.[] | "--- COMMENT \(.id) by \(.user.login) at \(.created_at) ---\n\(.body)\n"'
```

Extract from the ticket:
- **Cluster ID** (e.g. `d6e181rl0lej9bi18tdg`) — appears in `ic oc cluster get` output, victbot analysis, or SRE notes
- **Cluster name** (e.g. `scp-2d70915d69e1`)
- **Carrier** (`ActualDatacenterCluster`, e.g. `prod-lon04.carrier106`)
- **Region** (e.g. `eu-gb` → maps to `uk-south` for xo commands)
- **Type**: Satellite cruiser (`multishift_cruiser`) vs classic IKS/ROKS

### 2. High-level status via Slack xo bot

Run these in parallel with `Promise.all` in `execute`:

```js
// Basic cluster info
tools.slack.slack_command({ command: "cluster <clusterID>" })

// Open health issues
tools.slack.slack_command({ command: "health.issues cluster=<clusterID>" })
```

Key fields to read from `cluster` output:
- `HealthState` / `HealthStatus` — `normal` is good
- `ActualState` / `DesiredState` — both should be `deployed`
- Worker counts by actual/health state
- `ActualDatacenterCluster` — the carrier this cruiser lives on

For more detail add `show=all` to the cluster command.

### 3. Deeper inspection via armada kubectl

**Important:** `armada_kubectl` needs the carrier/tugboat's **opaque cluster ID**,
not the human-readable name like `prod-lon04-carrier106`. To find it:

```js
// Try the SQL query (table is `cluster`, not `clusters`)
tools.slack.slack_command({
  command: "sqlQuery region=<xo-region> SELECT cluster_id, name FROM cluster WHERE name = '<carrier-name>' LIMIT 3"
})
```

If that returns no results (carrier names aren't always in the DB), get it from
a Jenkins GMI job log:

```bash
curl -sS -u "${JENKINS_USER}:${JENKINS_API_KEY}" \
  "https://alchemy-containers-jenkins.swg-devops.com/job/Containers-Runtime/job/armada-deploy-get-master-info/<job-N>/consoleText" \
  | grep -E "datacenter_cluster|responding_carrier|TUGBOAT_CLUSTER"
```

Once you have the carrier cluster ID, run kubectl against it. For a **Satellite
cruiser**, the control plane lives in a namespace named after the cluster ID:

```js
tools.armada.armada_kubectl({
  cluster_id: "<carrier-opaque-id>",
  commands: [
    "kubectl get pods -n <cruiserClusterID> -o wide",
    "kubectl get etcdcluster -n <cruiserClusterID>",
    "kubectl get nodes -o wide"
  ]
})
```

For a **classic IKS/ROKS** cluster (non-Satellite), pass the cluster's own ID directly.

### 4. Region → xo region mapping

| IBM Cloud region | xo region name |
|---|---|
| `eu-gb` | `uk-south` |
| `eu-de` | `eu-central` |
| `us-south` | `us-south` |
| `us-east` | `us-east` |
| `au-syd` | `ap-south` |
| `jp-tok` | `ap-north` |
| `ca-tor` | `ca-tor` |
| `br-sao` | `br-sao` |

### 5. Analyse and report

Structure the report as:

1. **Cluster summary** — ID, name, type, carrier, region
2. **Current health** — HealthState, worker counts, any open issues
3. **Root cause** — what failed and why (timeline if relevant)
4. **Evidence** — key pod states, etcd status, error messages
5. **Resolution** — what was/needs to be done, data loss window if applicable
6. **Customer action items** — anything the customer must do post-recovery

## Common failure patterns for Satellite clusters

### etcd quorum loss
- **Symptom:** `HealthState: error`, master pods in `CrashLoopBackOff` / `Init:0/N`, DNS has zero endpoints
- **Cause:** All etcd-hosting nodes replaced/rebooted simultaneously (quorum needs N/2+1 members)
- **Fix:** `armada-restore-cluster-etcd` Jenkins job; restore from latest backup
- **Data loss:** window between backup timestamp and outage start
- **Post-restore:** customer must reboot workers, re-delete/reload workers added/deleted during loss window, reapply app changes

### kube-apiserver unreachable
- **Symptom:** `EOF` or connection refused on cluster API; workers show `normal` but master `error`
- **Check:** etcd health first (often the root cause), then kube-apiserver pod logs

### Location Action Required (R0025 / R0002)
- **Symptom:** `Status: Location Action Required` on the cluster
- **Meaning:** Satellite location has hosts with issues; check the location's host health
- **Check:** `health.issues cluster=<id>` and look at Satellite location status

## Jenkins credentials

```bash
JENKINS_USER   # set as env var (e.g. tamas.mojzes@ibm.com)
JENKINS_API_KEY  # set as env var, also in ~/keys/jenkins-api-key
```

GMI job: `https://alchemy-containers-jenkins.swg-devops.com/job/Containers-Runtime/job/armada-deploy-get-master-info/`
etcd restore: `https://alchemy-containers-jenkins.swg-devops.com/job/Containers-Runtime/job/armada-restore-cluster-etcd/`

## Notes

- The victbot comment (posted automatically) contains a detailed health check — always read it first before running your own queries; it may already have the answer.
- The SRE notes comment (usually from `Laszlo-Vigh` or similar) contains the diagnosis.
- The dev notes comment contains the remediation action and job links.
- `github.ibm.com` requires `gh` CLI; direct `webfetch` will be redirected to SSO login.
- Slack xo `sqlQuery` table is `cluster` (not `clusters`); columns vary — use `SELECT *` with `LIMIT 1` to discover schema if needed.
