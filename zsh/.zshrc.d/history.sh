# History settings
export HISTFILE="$HOME/.zsh_eternal_history"
export HISTSIZE=1000000
export SAVEHIST=1000000

# Append to history file immediately after each command, share across sessions.
# EXTENDED_HISTORY prefixes each entry with ": <timestamp>:<elapsed>;" natively.
setopt INC_APPEND_HISTORY
setopt SHARE_HISTORY
setopt EXTENDED_HISTORY

# Don't store duplicate entries or commands starting with a space.
setopt HIST_IGNORE_DUPS
setopt HIST_IGNORE_SPACE
