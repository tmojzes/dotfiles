# Enables starship.
if command -v starship &>/dev/null; then
    eval "$(starship init bash)"
fi

if command -v direnv &>/dev/null; then
    eval "$(direnv hook bash)"
fi

# Vite+ bin (https://viteplus.dev)
if [ -d "$HOME/.vite-plus" ]; then
    . "$HOME/.vite-plus/env"
fi

# Ask which tmux session to attach to or create a new session.
if command -v tmux &>/dev/null && [ -z "$TMUX" ]; then
    echo "Current tmux sessions:"
    tmux ls
    read -p "Enter session name to attach, or press enter for new: " session_name
    if [ -z "$session_name" ]; then
        tmux new-session
    else
        tmux attach-session -t "$session_name"
    fi
fi
