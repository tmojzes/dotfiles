# Enables starship.
if command -v starship &>/dev/null; then
    eval "$(starship init bash)"
fi

# Starts tmux session.
if command -v tmux &>/dev/null; then
    # Start tmux session if shell is not running in one already.
    [[ -z "$TMUX" ]] && exec tmux
fi
