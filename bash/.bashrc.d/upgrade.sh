upgrade() {
    # Helpers for colored output
    print_step() { echo -e "\n\033[1;34m=== $1 ===\033[0m"; }
    print_title() { echo -e "\n\033[1;32m--- $1 ---\033[0m"; }

    print_title "Starting System Updates"

    print_step "System Packages"
    if command -v apt &>/dev/null; then
        sudo apt update && sudo apt upgrade -y
    elif command -v dnf &>/dev/null; then
        sudo dnf upgrade -y
    elif command -v paru &>/dev/null; then
        paru -Syyu --noconfirm
    elif command -v zypper &>/dev/null; then
        sudo zypper update -y
    fi

    # 2. Reusable helper for optional tools
    update_tool() {
        local cmd=$1
        local msg=$2

        shift 2
        if command -v "$cmd" &>/dev/null; then
            print_step "$msg"
            "$@"
        fi
    }

    update_tool snap "Updating Snaps" sudo snap refresh
    update_tool flatpak "Updating Flatpaks" flatpak update -y
    update_tool go-global-update "Updating Go packages" go-global-update
    update_tool rustup "Updating Rust toolchain" rustup update

    if command -v cargo &>/dev/null; then
        print_step "Updating Rust crates"
        if ! cargo install-update -V &>/dev/null; then
            cargo install cargo-update
        fi
        cargo install-update -a
    fi

    if command -v ibmcloud &>/dev/null && [ -d "$HOME/ibmcloud_homes" ]; then
        print_step "Updating IBM Cloud Plugins"

        for ic_home in "$HOME/ibmcloud_homes"/*; do
            if [ -d "$ic_home" ]; then
                echo "-> Profile: $(basename "$ic_home")"
                IBMCLOUD_HOME="$ic_home" ibmcloud plugin update --all
            fi
        done
    fi

    print_title "All updates complete!"
}
