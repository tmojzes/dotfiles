# Use modern completion system
autoload -Uz compinit
compinit

autoload -Uz promptinit
promptinit
prompt adam1

# Keep 10000 lines of history within the shell and save it to ~/.zsh_history:
HISTSIZE=10000
SAVEHIST=10000
HISTFILE=~/.zsh_history

## Custom path config
# GOPATH for microservices testing
# export GOPATH=/home/tmojzes/Projects/microservice/testing_ms
ZIG_PATH="/usr/local/zig"
DOOM_EMACS_PATH="$HOME/.config/emacs/bin"
export PATH=$PATH:/usr/local/go/bin:$HOME/.cargo/bin:"${KREW_ROOT:-$HOME/.krew}/bin":$ZIG_PATH:$DOOM_EMACS_PATH

# Go env variables
export GONOPROXY="gitlabe2.ext.net.nokia.com"
export GONOSUMDB="gitlabe2.ext.net.nokia.com"
export GOPRIVATE="gitlabe2.ext.net.nokia.com"
export GO111MODULE="on"

# Completions
source <(talosctl completion zsh)
source <(kubectl completion zsh)
source <(operator-sdk completion zsh)
source <(kubectl krew completion zsh)
source <(k3d completion zsh)
source <(clusterctl completion zsh)
source <(helm completion zsh)

# Aliases
alias sopse="sops --encrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"
alias sopsd="sops --decrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"
alias upgrade="sudo apt update && sudo apt upgrade -y && flatpak update -y && go-global-update"
alias emacs="emacsclient -c -a 'emacs'"


bindkey -v
export EDITOR="lvim"
export SOPS_AGE_KEY_FILE=$HOME/.sops/key

eval "$(starship init zsh)"
autoload -U +X bashcompinit && bashcompinit
complete -o nospace -C /home/tmojzes/go/bin/gocomplete go
