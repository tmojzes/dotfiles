#
# ~/.bashrc
#

# If not running interactively, don't do anything
[[ $- != *i* ]] && return

# Set aliases
alias ls='ls --color=auto'
alias ll='ls -lah --color=auto'
alias grep='grep --color=auto'
alias vim="nvim"

# alias sopse="sops --encrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"
# alias sopsd="sops --decrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"

# Set exports
export EDITOR="nvim"
export VISUAL="nvim"
export PATH=$PATH:$HOME/go/bin:$HOME/.local/bin:$HOME/.cargo/bin
export XDG_DATA_DIRS=$HOME/.local/share/flatpak/exports/share:/var/lib/flatpak/exports/share:/usr/local/share:/usr/share
export XDG_CONFIG_HOME=$HOME/.config
export XDG_DATA_HOME=$HOME/.local/share/

# export SOPS_AGE_KEY_FILE=$HOME/.local/share/sops/key

# Completions
# source <(talosctl completion bash)
# source <(kubectl completion bash)
# source <(operator-sdk completion bash)
# source <(kustomize completion bash)
# source <(k3d completion bash)
# source <(clusterctl completion bash)
# source <(helm completion bash)
# source <(pdm completion bash)
# complete -C /home/tmojzes/go/bin/gocomplete go

# Set vi mode
set -o vi

eval "$(starship init bash)"
