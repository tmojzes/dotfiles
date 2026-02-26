upgrade() {
    echo "--- Starting System Updates ---"

    # 1. System Package Manager Detection
    if command -v apt &>/dev/null; then
        sudo apt update && sudo apt upgrade -y
    elif command -v dnf &>/dev/null; then
        sudo dnf upgrade -y
    elif command -v paru &>/dev/null; then
        sudo paru -Syyu --noconfirm
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

    # Rust update
    if command -v rustup &>/dev/null; then
        echo "Updating Rust toolchain..."
        rustup update
    fi

    # Update installed crates
    if command -v cargo &>/dev/null; then
        echo "Updating Rust crates..."
        cargo install cargo-update
        cargo install-update -a
    fi

    echo "--- All updates complete! ---"
}
