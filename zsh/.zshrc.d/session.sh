# Enables starship.
if command -v starship &>/dev/null; then
    eval "$(starship init zsh)"
fi

# Starts tmux session.
if command -v tmux &>/dev/null; then
    # Start tmux session if shell is not running in one already.
    [[ -z "$TMUX" ]] && exec tmux
fi

# Direnv hook setup
if command -v direnv &>/dev/null; then
    eval "$(direnv hook zsh)"
fi
