# Enables starship.
if command -v starship &>/dev/null; then
    eval "$(starship init zsh)"
fi

# Vite+ bin (https://viteplus.dev)
if [ -d "$HOME/.vite-plus" ]; then
    . "$HOME/.vite-plus/env"
fi

# Starts tmux session (Linux only).
if command -v tmux &>/dev/null && [[ "$(uname -s)" == "Linux" ]]; then
    # Start tmux session if shell is not running in one already.
    [[ -z "$TMUX" ]] && exec tmux
fi

# Direnv hook setup
if command -v direnv &>/dev/null; then
    eval "$(direnv hook zsh)"
fi
