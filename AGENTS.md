# AGENTS.md

Personal GNU stow dotfiles repo. No build, tests, or CI — changes are verified by using the live system.

## Layout & deployment

- Each top-level directory is a stow package mirroring `$HOME` (e.g. `nvim/.config/nvim/x` → `~/.config/nvim/x`). Create/edit files inside the package; live configs are symlinks back into this repo.
- Deploy from repo root: `stow <pkg>` (or `stow */` for all; the trailing `/` globs directories only). Preview first with `stow -n -v <pkg>`; unlink with `stow -D <pkg>`.
- Not every package is stowed on every machine — a missing symlink in `$HOME` doesn't mean a package is unused.
- `scripts/lf/` is a nested package: deploy with `stow -d scripts lf`. Plain `stow scripts` would create a stray `~/lf`.
- `.stow-local-ignore` already excludes README/LICENSE from stowing.

## Generated files inside the working tree

Configs are symlinked at directory level, so tools write generated content into the repo. These are gitignored — never commit or delete them:
- `nvim/.config/nvim/lazy-lock.json`, `nvim/.config/nvim/lazyvim.json` (lazy.nvim/LazyVim)
- `tmux/.config/tmux/plugins/` (TPM)

## Migrated Packages
The following configurations and packages were migrated to personal NixOS/Home-Manager infrastructure in `infra`:
- `bob`, `fonts`, `gemini`, `ghostty`, `git`, `lazygit`, `nix` (Home Manager flake), `opencode`, and `shell`.
They are managed declaratively in `~/infra/modules/home/`.

## Package notes

- **nvim**: LazyVim distro config. Prefer adding a `lazyvim.plugins.extras.*` import in `lua/config/lazy.lua` over hand-written specs in `lua/plugins/` (recent commits migrate toward extras). Lua style: tabs, width 4, 100 columns (`stylua.toml`).
- **mcp**: Go MCP servers for opencode. Source lives at `mcp/mcp/` (package `mcp` stows `~/mcp → dotfiles/mcp/mcp` — the doubled name is the `opencode/.config/opencode/` pattern, don't flatten it or stow scatters files into `~/`). Single Go module, one dir per server under `cmd/<name>/`; `go build` output goes to `~/.local/bin/<name>-mcp`, which is a stow fold over `scripts/.local/bin/` — binaries land there as gitignored artifacts (root `.gitignore`, never commit them). All tooling (gofumpt, staticcheck, govulncheck — and go-task itself) is pinned via `go.mod` `tool` directives and invoked as `go tool <name>`; nix's `task`/`go` are only bootstraps. Run tasks from the module dir: `cd ~/mcp && go tool task build|test|fmt|lint|vuln` (first `go tool task` call compiles task into the module cache once). Current servers: `armada` (readonly kubectl via a Jenkins job; `JENKINS_USER`/`JENKINS_API_KEY`) and `slack` (armada-xo bot commands posted as replies under the latest ':thread:debug satellite' thread in #armada-xo on ibm.enterprise.slack.com; `SLACK_XOXC`+`SLACK_XOXD` session credentials, both required, expire with the Slack session — refresh with the slack-token skill). New servers: add `cmd/<name>/main.go`, register it under `mcp.servers` in the opencode package's global `opencode.json`, rebuild, and restart the service. Credentials are hand-exported from `~/keys` (like `BOB_API_KEY`), never committed.
- **pacman**: contains only `makepkg.conf`.

## Known stale/broken

- README and `scripts/.local/bin/{install,list}-*.sh` reference a `programs/` package-list directory that no longer exists — those scripts are broken; don't rely on them or "fix" them against that path.

## Git workflow

- Conventional commits (`feat(nvim): ...`, `fix: ...`), committed directly to `main`.
