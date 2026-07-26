{ pkgs, ... }:

{
  home.username = "tmojzes";
  # home.homeDirectory is set per-host in flake.nix.
  home.stateVersion = "26.05";

  # Determinate Nix manages /etc/nix/nix.conf; don't let home-manager touch it.
  nix.enable = false;

  programs.home-manager.enable = true;

  home.packages = with pkgs; [
    # Editors & terminal
    helix
    neovim
    tmux
    zellij

    # Languages & runtimes (rust stays on rustup, odin on brew)
    elixir
    go
    nodejs
    bun
    python313
    tinygo
    zig
    luarocks

    # Build & task runners
    go-task
    just

    # Git & GitHub
    gh
    git-lfs
    jujutsu
    lazygit

    # Linters & formatters
    ast-grep
    black
    go-tools
    golines
    golangci-lint
    isort
    markdownlint-cli
    prettier
    pylint
    shellcheck
    yamllint

    # Python tooling (pipx and poetry stay on brew)
    uv

    # CLI utilities
    direnv
    fd
    fzf
    ripgrep
    tree-sitter
    yq-go

    # Shell & dotfiles tooling
    starship
    stow

    # Docs & diagrams
    markdown-toc
    mermaid-cli
    tectonic

    # Web
    tailwindcss

    # AI tools
    codex
    gemini-cli
    goose-cli
    ollama
    opencode

    # Infra & Kubernetes
    age
    sops
    cilium-cli
    k9s
    ko
    kubebuilder
    # minikube also ships a kubectl binary; prefer the standalone one.
    (lib.hiPrio kubectl)
    kustomize
    minikube
    molecule
    openshift
    opentofu
    operator-sdk
    talhelper
    talosctl
    terraform

    # Virtualization
    qemu
    vfkit
  ];
}
