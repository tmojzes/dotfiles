# Set exports
export EDITOR="nvim"
export VISUAL="nvim"
export FLYCTL_INSTALL=${HOME}/.fly
export XDG_DATA_DIRS=${HOME}/.local/share/flatpak/exports/share:/var/lib/flatpak/exports/share:/usr/local/share:/usr/share
export XDG_CONFIG_HOME=${HOME}/.config
export XDG_DATA_HOME=${HOME}/.local/share/
export XDG_RUNTIME_DIR=/run/user/$(id -u)
export WAYLAND_DISPLAY=wayland-1
export GOBIN="${HOME}/go/bin/"
export ODIN_BIN="${HOME}/projects/oss/Odin/"
export BUN_BIN="${HOME}/.bun/bin"

export PATH=${PATH}:${HOME}/go/bin:${HOME}/.local/bin:${HOME}/.cargo/bin:${FLYCTL_INSTALL}/bin:${ODIN_BIN}:${BUN_BIN}
