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

if command -v gocomplete &>/dev/null; then
    complete -C "$HOME/go/bin/gocomplete" go
fi

if command -v rustup &>/dev/null; then
    source <(rustup completions bash)
fi

if command -v tailscale &>/dev/null; then
    source <(tailscale completion bash)
fi

if command -v minikube &>/dev/null; then
    source <(minikube completion bash)
fi
