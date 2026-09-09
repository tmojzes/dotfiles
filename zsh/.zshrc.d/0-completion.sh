# load completion
autoload -U +X compinit && compinit

# Completions
if (("$+commands[talosctl]")); then
    source <(talosctl completion zsh)
fi

if (("$+commands[kubectl]")); then
    source <(kubectl completion zsh)
fi

if (("$+commands[operator-sdk]")); then
    source <(operator-sdk completion zsh)
fi

if (("$+commands[kustomize]")); then
    source <(kustomize completion zsh)
fi

if (("$+commands[k3d]")); then
    source <(k3d completion zsh)
fi

if (("$+commands[clusterctl]")); then
    source <(clusterctl completion zsh)
fi

if (("$+commands[helm]")); then
    source <(helm completion zsh)
fi

if (("$+commands[gocomplete]")); then
    complete -C "$HOME/go/bin/gocomplete" go
fi

if (("$+commands[rustup]")); then
    source <(rustup completions zsh)
fi

if (("$+commands[tailscale]")); then
    source <(tailscale completion zsh)
fi

if (("$+commands[ibmcloud]")); then
    source /usr/local/ibmcloud/autocomplete/zsh_autocomplete
fi
