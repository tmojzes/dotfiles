{ pkgs, ... }:

{
  home.username = "tmojzes";
  # home.homeDirectory is set per-host in flake.nix.
  home.stateVersion = "26.05";

  # Determinate Nix manages /etc/nix/nix.conf; don't let home-manager touch it.
  nix.enable = false;

  programs.home-manager.enable = true;

  home.packages = with pkgs; [
    go
    nodejs
    bun
    lazygit
    neovim
    ripgrep
    tinygo
    tmux
  ];
}
