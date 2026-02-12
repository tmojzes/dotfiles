upgrade() {
    echo "--- Starting System Updates ---"

    # 1. System Package Manager Detection
    if command -v apt &>/dev/null; then
        sudo apt update && sudo apt upgrade -y
    elif command -v dnf &>/dev/null; then
        sudo dnf upgrade -y
    elif command -v pacman &>/dev/null; then
        sudo pacman -Syyu --noconfirm
    elif command -v zypper &>/dev/null; then
        sudo zypper update -y
    fi

    # Snap Updates
    if command -v snap &>/dev/null; then
        echo "Updating Snaps..."
        sudo snap refresh
    fi

    # Flatpak Updates
    if command -v flatpak &>/dev/null; then
        echo "Updating Flatpaks..."
        flatpak update -y
    fi

    # Go Global Update
    if command -v go-global-update &>/dev/null; then
        echo "Updating Go packages..."
        go-global-update
    fi

    echo "--- All updates complete! ---"
}
