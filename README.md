# Dotfiles

![dotfiles image](./dotfiles.png)

## Installing

You will need `git` and GNU `stow`

Clone into your `$HOME` directory or `~`

```bash
git clone https://github.com/tmojzes/dotfiles.git ~
```

Run `stow` to symlink everything or just select what you want

```bash
stow */ # Everything (the '/' ignores the README)
```

```bash
stow nvim # Just my neovim config
```

## Dev environment (Nix)

Dev tools (languages, editors, linters, git/CLI utilities, AI CLIs, k8s/infra tooling) are
managed by home-manager via the `nix` package. Graphical applications are managed
declaratively per platform: on Linux through `nix-flatpak` (`gui.nix`), and on macOS
through `home-manager-brew` (`gui-darwin.nix`) for casks; no brew formulae are installed.

1. Install [Determinate Nix](https://determinate.systems/nix/):

   ```bash
   curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install
   ```

2. Stow the config and activate it:

   ```bash
   stow nix
   # Bootstrap with home-manager (pick the target matching your host architecture):
   nix run github:nix-community/home-manager -- switch --flake "$HOME/.config/home-manager#tmojzes-x86_64-linux" # x86_64 Linux
   nix run github:nix-community/home-manager -- switch --flake "$HOME/.config/home-manager#tmojzes"              # aarch64 Linux
   nix run github:nix-community/home-manager -- switch --flake "$HOME/.config/home-manager#tmojzes-mac"          # macOS
   ```

After the first activation, use `home-manager switch` (or the `upgrade` function) instead of `nix run ...`.

## Programs

An updated list of all the programs I use can be found in the `programs` directory
