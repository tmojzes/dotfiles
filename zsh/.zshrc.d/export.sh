# Set exports
export EDITOR="nvim"
export VISUAL="nvim"
export FLYCTL_INSTALL="$HOME/.fly"
export XDG_DATA_DIRS=$HOME/.local/share/flatpak/exports/share:/var/lib/flatpak/exports/share:/usr/local/share:/usr/share
export XDG_CONFIG_HOME=$HOME/.config
export XDG_DATA_HOME=$HOME/.local/share/
export XDG_RUNTIME_DIR=/run/user/$(id -u)
export WAYLAND_DISPLAY=wayland-0
export GOBIN="$HOME/go/bin/"
export GO_INSTALL_PATH="/usr/local/go/bin"
export ODIN_BIN="$HOME/projects/oss/Odin/"
export BUN_BIN="$HOME/.bun/bin"
export CARGO_BIN="$HOME/.cargo/bin"

[[ -d "$HOME/.rd/bin" ]] && export RANCHER_BIN="$HOME/.rd/bin"

export PATH=$PATH:$GOBIN:$HOME/.local/bin:$CARGO_BIN:$FLYCTL_INSTALL/bin:$ODIN_BIN:$BUN_BIN:${RANCHER_BIN:+$RANCHER_BIN:}$GO_INSTALL_PATH
