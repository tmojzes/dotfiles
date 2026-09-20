# Set exports
export EDITOR="nvim"
export VISUAL="nvim"
export FLYCTL_INSTALL="$HOME/.fly"
export XDG_DATA_DIRS="$HOME/.local/share/flatpak/exports/share:/var/lib/flatpak/exports/share:/usr/local/share:/usr/share"
export XDG_CONFIG_HOME="$HOME/.config"
export XDG_DATA_HOME="$HOME/.local/share/"
export XDG_RUNTIME_DIR="/run/user/$(id -u)"
if command -v systemctl &>/dev/null; then
    export WAYLAND_DISPLAY=$(systemctl --user show-environment | grep '^WAYLAND_DISPLAY=' | cut -d= -f2)
fi
export GOBIN="$HOME/go/bin/"
export NPM_BIN="$(npm config get prefix)/bin"
export GO_INSTALL_PATH="/usr/local/go/bin"
export ODIN_BIN="$HOME/projects/oss/Odin/"
export BUN_BIN="$HOME/.bun/bin"
export CARGO_BIN="$HOME/.cargo/bin"
export OPENCODE_BIN="$HOME/.opencode/bin"
export NODE_USER_BIN="$HOME/.npm-global/bin"

export PATH="$PATH:$GOBIN:$HOME/.local/bin:$HOME/.cargo/bin:$FLYCTL_INSTALL/bin:$ODIN_BIN:$BUN_BIN:${RANCHER_BIN:+$RANCHER_BIN:}$GO_INSTALL_PATH:$OPENCODE_BIN:$NPM_BIN:$NODE_USER_BIN"
