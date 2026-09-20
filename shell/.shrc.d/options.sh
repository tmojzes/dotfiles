# Set vi mode (works in both bash and zsh)
set -o vi

# Zsh-specific keybindings
if [ -n "$ZSH_VERSION" ]; then
    # Set backward delete with backspace
    bindkey -v '^?' backward-delete-char
    # Set ctrl + r for search in zsh history
    bindkey -v '^R' history-incremental-search-backward
fi
