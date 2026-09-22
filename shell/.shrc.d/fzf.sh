# Fuzzy finder (fzf) integration
if command -v fzf &>/dev/null; then
    if [ -n "$ZSH_VERSION" ]; then
        eval "$(fzf --zsh)"
    elif [ -n "$BASH_VERSION" ]; then
        eval "$(fzf --bash)"
    fi
fi
