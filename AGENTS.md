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

Exception: `nix/.config/home-manager/flake.lock` IS committed on purpose (reproducible package set) — don't gitignore or delete it.

## Package notes

- **nvim**: LazyVim distro config. Prefer adding a `lazyvim.plugins.extras.*` import in `lua/config/lazy.lua` over hand-written specs in `lua/plugins/` (recent commits migrate toward extras). Lua style: tabs, width 4, 100 columns (`stylua.toml`).
- **shell**: merged bash+zsh config (replaced the former `bash/` and `zsh/` packages). `.bashrc` and `.zshrc` are thin wrappers that source every file in `~/.shrc.d/`; add new shell behavior as a new drop-in file there instead of editing the rc files. Drop-ins must work in both shells: guard shell-specific code with `if [ -n "$ZSH_VERSION" ]` / `elif [ -n "$BASH_VERSION" ]`, detect tools with `command -v` (not zsh-only `$+commands`), use POSIX function syntax. `00-completion.sh` must keep its numeric prefix so zsh's `compinit`/`bashcompinit` run before `compdef` in `alias.sh`. Bash completion data files live in `shell/.local/share/bash-completion/completions/`. Migrating a machine to this package: `stow -D bash zsh` (whichever are stowed there), then `stow shell`.
- **git**: `.gitconfig` always includes `~/.gitconfig_base` and conditionally `~/.gitconfig_ibm` (`includeIf "gitdir/i:~/projects/ibm/"`). Work-specific identity/settings go in `.gitconfig_ibm`, personal in `.gitconfig_base`.
- **opencode**: the user's *global* opencode config deployed via stow — not repo-local config for working in this repo.
- **mcp**: Go MCP servers for opencode. Source lives at `mcp/mcp/` (package `mcp` stows `~/mcp → dotfiles/mcp/mcp` — the doubled name is the `opencode/.config/opencode/` pattern, don't flatten it or stow scatters files into `~/`). Single Go module, one dir per server under `cmd/<name>/`; `go build` output goes to `~/.local/bin/<name>-mcp`, which is a stow fold over `scripts/.local/bin/` — binaries land there as gitignored artifacts (root `.gitignore`, never commit them). All tooling (gofumpt, staticcheck, govulncheck — and go-task itself) is pinned via `go.mod` `tool` directives and invoked as `go tool <name>`; nix's `task`/`go` are only bootstraps. Run tasks from the module dir: `cd ~/mcp && go tool task build|test|fmt|lint|vuln` (first `go tool task` call compiles task into the module cache once). Current servers: `armada` (readonly kubectl via a Jenkins job; `JENKINS_USER`/`JENKINS_API_KEY`) and `slack` (armada-xo bot commands posted as replies under the latest ':thread:debug satellite' thread in #armada-xo on ibm.enterprise.slack.com; `SLACK_XOXC`+`SLACK_XOXD` session credentials, both required, expire with the Slack session — refresh with the slack-token skill). New servers: add `cmd/<name>/main.go`, register it under `mcp.servers` in the opencode package's global `opencode.json`, rebuild, and restart the service. Credentials are hand-exported from `~/keys` (like `BOB_API_KEY`), never committed.
- **pacman**: contains only `makepkg.conf`.
- **nix**: home-manager flake (packages-only — dotfiles stay stow-managed). Dev tools live in `home.nix`: editors/terminal (helix, neovim, tmux, zellij), languages (elixir, go, nodejs, bun, python313, tinygo, zig, luarocks), linters/formatters, git tools (gh, git-lfs, jujutsu, lazygit), CLI utils (direnv, fd, fzf, ripgrep, tree-sitter, yq-go), shell/dotfiles tooling (starship, stow), AI CLIs (codex, gemini-cli, goose-cli, ollama, opencode), and infra/k8s tooling (sops, k9s, kubectl, kustomize, opentofu, talosctl, etc.). Graphical applications are managed declaratively per platform: on Linux via `nix-flatpak` in `gui.nix` (Chrome, Spotify, Vesktop, Slack, Bitwarden), and on macOS via `home-manager-brew` in `gui-darwin.nix` (Chrome, Spotify, Slack, Discord, Ghostty, Bitwarden casks). rustup comes from nixpkgs but still manages toolchains in ~/.rustup (rust stays rustup-managed, never nix-built). Homebrew is casks-only now (GUI apps + nerd fonts); all formulae were removed. Lost with the formulae, reinstall externally if needed: pipx + poetry (user preference), odin (LLVM 18 build broken against macOS 26 SDK), dagger/cpm (not in nixpkgs). staticcheck comes from nixpkgs `go-tools` (which is honnef.co/go/tools, not golang.org/x/tools). flake.nix allows unfree only for `terraform` (`allowUnfreePredicate`). Three hosts in `flake.nix`: `tmojzes` (aarch64-linux, the bare `--flake` default), `tmojzes-x86_64-linux` (x86_64-linux), and `tmojzes-mac` (aarch64-darwin); `home.homeDirectory` is set per-host in the flake. Deploy: `stow nix` then `home-manager switch --flake ~/.config/home-manager#tmojzes` (aarch64 Linux), `...#tmojzes-x86_64-linux` (x86_64 Linux), or `...#tmojzes-mac` (macOS). `upgrade.sh` auto-selects the right host from `uname`/`uname -m`. Bootstrap on a fresh machine: install Determinate Nix (`curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install`), then `nix run github:nix-community/home-manager -- switch --flake ~/.config/home-manager#<host>` (`tmojzes-x86_64-linux`, `tmojzes`, or `tmojzes-mac`). Gotchas: flakes only see git-tracked files, so `git add` new files in this package before switching; `nix.enable = false` because Determinate Nix owns `/etc/nix/nix.conf`; Nix itself self-upgrades via `determinate-nixd` (not in `upgrade.sh`).

## Known stale/broken

- README and `scripts/.local/bin/{install,list}-*.sh` reference a `programs/` package-list directory that no longer exists — those scripts are broken; don't rely on them or "fix" them against that path.

## Git workflow

- Conventional commits (`feat(nvim): ...`, `fix: ...`), committed directly to `main`.
