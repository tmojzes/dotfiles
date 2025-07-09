# load completion
autoload -U +X compinit && compinit

# Completions
if command -v talosctl &>/dev/null; then
    source <(talosctl completion zsh)
fi

if command -v kubectl &>/dev/null; then
    source <(kubectl completion zsh)
fi

if command -v operator-sdk &>/dev/null; then
    source <(operator-sdk completion zsh)
fi

if command -v kustomize &>/dev/null; then
    source <(kustomize completion zsh)
fi

if command -v k3d &>/dev/null; then
    source <(k3d completion zsh)
fi

if command -v clusterctl &>/dev/null; then
    source <(clusterctl completion zsh)
fi

if command -v helm &>/dev/null; then
    source <(helm completion zsh)
fi

if command -v gocomplete &>/dev/null; then
    complete -C "$HOME/go/bin/gocomplete" go
fi

if command -v rustup &>/dev/null; then
    source <(rustup completions zsh)
fi

if command -v tailscale &>/dev/null; then
    source <(tailscale completion zsh)
fi

if command -v ibmcloud &>/dev/null; then
    source /usr/local/ibmcloud/autocomplete/zsh_autocomplete
fi
