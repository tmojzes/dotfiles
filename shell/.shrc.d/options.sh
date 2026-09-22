# Set vi mode (works in both bash and zsh)
set -o vi

# Zsh-specific keybindings
if [ -n "$ZSH_VERSION" ]; then
    # Set backward delete with backspace
    bindkey -v '^?' backward-delete-char
    # Fallback ctrl + r for search in zsh history when fzf is not installed
    if ! command -v fzf &>/dev/null; then
        bindkey -v '^R' history-incremental-search-backward
    fi
fi
