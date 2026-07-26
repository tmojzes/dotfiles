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
- **bash / zsh**: `.bashrc` and `.zshrc` source every file in `~/.bashrc.d/` / `~/.zshrc.d/`. Add new shell behavior as a new drop-in file there instead of editing the rc files.
- **git**: `.gitconfig` always includes `~/.gitconfig_base` and conditionally `~/.gitconfig_ibm` (`includeIf "gitdir/i:~/projects/ibm/"`). Work-specific identity/settings go in `.gitconfig_ibm`, personal in `.gitconfig_base`.
- **opencode**: the user's *global* opencode config deployed via stow — not repo-local config for working in this repo.
- **pacman**: contains only `makepkg.conf`.
- **nix**: home-manager flake (packages-only — dotfiles stay stow-managed). Dev tools (go, nodejs, bun, lazygit, neovim, ripgrep, tinygo, tmux) live in `home.nix`; rust stays on rustup. Two hosts in `flake.nix`: `tmojzes` (aarch64-linux, the bare `--flake` default) and `tmojzes-mac` (aarch64-darwin); `home.homeDirectory` is set per-host in the flake. Deploy: `stow nix` then `home-manager switch --flake ~/.config/home-manager#tmojzes` (Linux) or `...#tmojzes-mac` (macOS). Bootstrap on a fresh machine: install Determinate Nix (`curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install`), then `nix run home-manager/master -- switch --flake ~/.config/home-manager#tmojzes-mac`. Gotchas: flakes only see git-tracked files, so `git add` new files in this package before switching; `nix.enable = false` because Determinate Nix owns `/etc/nix/nix.conf`; Nix itself self-upgrades via `determinate-nixd` (not in `upgrade.sh`).

## Known stale/broken

- README and `scripts/.local/bin/{install,list}-*.sh` reference a `programs/` package-list directory that no longer exists — those scripts are broken; don't rely on them or "fix" them against that path.

## Git workflow

- Conventional commits (`feat(nvim): ...`, `fix: ...`), committed directly to `main`.
