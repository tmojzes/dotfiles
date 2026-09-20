# Enables starship.
if command -v starship &>/dev/null; then
    if [ -n "$ZSH_VERSION" ]; then
        eval "$(starship init zsh)"
    else
        eval "$(starship init bash)"
    fi
fi

if command -v direnv &>/dev/null; then
    if [ -n "$ZSH_VERSION" ]; then
        eval "$(direnv hook zsh)"
    else
        eval "$(direnv hook bash)"
    fi
fi

# Vite+ bin (https://viteplus.dev)
if [ -d "$HOME/.vite-plus" ]; then
    . "$HOME/.vite-plus/env"
fi

# Ask which tmux session to attach to or create a new session (Linux only).
if command -v tmux &>/dev/null && [ -z "$TMUX" ] && [ "$(uname -s)" = "Linux" ]; then
    echo "Current tmux sessions:"
    tmux ls
    printf "Enter session name to attach, or press enter for new: "
    read session_name
    if [ -z "$session_name" ]; then
        tmux new-session
    else
        tmux attach-session -t "$session_name"
    fi
fi
