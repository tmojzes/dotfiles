#
# ~/.bashrc
#

# If not running interactively, don't do anything
[[ $- != *i* ]] && return

if [ -f /etc/bash_completion ] && ! shopt -oq posix; then
    . /etc/bash_completion
fi

# Set aliases
alias ls='ls --color=auto'
alias ll='ls -lah --color=auto'
alias grep='grep --color=auto'
alias vim="nvim"
alias upgrade="sudo apt update && sudo apt upgrade -y && flatpak update -y && go-global-update"

# alias sopse="sops --encrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"
# alias sopsd="sops --decrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"

# Set exports
export EDITOR="nvim"
export VISUAL="nvim"
export FLYCTL_INSTALL="/home/tmojzes/.fly"
export PATH=$PATH:$HOME/go/bin:$HOME/.local/bin:$HOME/.cargo/bin:$FLYCTL_INSTALL/bin:/usr/local/bin/Odin
export XDG_DATA_DIRS=$HOME/.local/share/flatpak/exports/share:/var/lib/flatpak/exports/share:/usr/local/share:/usr/share
export XDG_CONFIG_HOME=$HOME/.config
export XDG_DATA_HOME=$HOME/.local/share/
export GOBIN="$HOME/go/bin/"

# export SOPS_AGE_KEY_FILE=$HOME/.local/share/sops/key

# Enables starship.
if command -v starship &>/dev/null; then
    eval "$(starship init bash)"
fi

# Completions
if command -v talosctl &>/dev/null; then
    source <(talosctl completion bash)
fi

if command -v kubectl &>/dev/null; then
    source <(kubectl completion bash)
fi

if command -v operator-sdk &>/dev/null; then
    source <(operator-sdk completion bash)
fi

if command -v kustomize &>/dev/null; then
    source <(kustomize completion bash)
fi

if command -v k3d &>/dev/null; then
    source <(k3d completion bash)
fi

if command -v clusterctl &>/dev/null; then
    source <(clusterctl completion bash)
fi

if command -v helm &>/dev/null; then
    source <(helm completion bash)
fi

if command -v go &>/dev/null; then
    complete -C /home/tmojzes/go/bin/gocomplete go
fi

# Set vi mode
set -o vi
