# Completions
# Load zsh's completion system (bash's is sourced from /etc/bash_completion in ~/.bashrc).
if [ -n "$ZSH_VERSION" ]; then
    autoload -U +X compinit && compinit
    # Provide bash's `complete` builtin in zsh so the `complete -C` lines work in both shells.
    autoload -U +X bashcompinit && bashcompinit
fi

# Dialect for tools with a `completion <shell>` subcommand
if [ -n "$ZSH_VERSION" ]; then
    _sh=zsh
else
    _sh=bash
fi

if command -v talosctl &>/dev/null; then
    source <(talosctl completion "$_sh")
fi

if command -v kubectl &>/dev/null; then
    source <(kubectl completion "$_sh")
fi

if command -v operator-sdk &>/dev/null; then
    source <(operator-sdk completion "$_sh")
fi

if command -v kustomize &>/dev/null; then
    source <(kustomize completion "$_sh")
fi

if command -v k3d &>/dev/null; then
    source <(k3d completion "$_sh")
fi

if command -v clusterctl &>/dev/null; then
    source <(clusterctl completion "$_sh")
fi

if command -v helm &>/dev/null; then
    source <(helm completion "$_sh")
fi

if command -v minikube &>/dev/null; then
    source <(minikube completion "$_sh")
fi

if command -v rustup &>/dev/null; then
    source <(rustup completions "$_sh")
fi

if command -v tailscale &>/dev/null; then
    source <(tailscale completion "$_sh")
fi

if command -v gocomplete &>/dev/null; then
    complete -C "$HOME/go/bin/gocomplete" go
fi

if command -v terraform &>/dev/null; then
    complete -C "$(command -v terraform)" terraform
fi

if command -v ibmcloud &>/dev/null; then
    if [ -n "$ZSH_VERSION" ]; then
        if [ -f /usr/local/ibmcloud/autocomplete/zsh_autocomplete ]; then
            source /usr/local/ibmcloud/autocomplete/zsh_autocomplete
        fi
    elif [ -f "$HOME/.local/share/bash-completion/completions/ibmcloud" ]; then
        source "$HOME/.local/share/bash-completion/completions/ibmcloud"
    fi
fi

unset _sh
