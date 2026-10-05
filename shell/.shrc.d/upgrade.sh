upgrade() {
    # Helpers for colored output
    print_step() { echo -e "\n\033[1;34m=== $1 ===\033[0m"; }
    print_title() { echo -e "\n\033[1;32m--- $1 ---\033[0m"; }

    # Reusable helper for optional tools
    update_tool() {
        local cmd=$1
        local msg=$2

        shift 2
        if command -v "$cmd" &>/dev/null; then
            print_step "$msg"
            "$@"
        fi
    }

    bob_update() {
        curl -fsSL https://bob.ibm.com/download/bobshell.sh | bash
    }

    print_title "Starting System Updates"

    if command -v bootc &>/dev/null; then
        print_step "System packages"
        sudo bootc update
    elif command -v apt &>/dev/null; then
        print_step "System Packages"
        sudo apt update && sudo apt upgrade -y
    elif command -v dnf &>/dev/null; then
        print_step "System Packages"
        sudo dnf upgrade -y
    elif command -v paru &>/dev/null; then
        print_step "System Packages"
        paru -Syyu --noconfirm
    elif command -v zypper &>/dev/null; then
        print_step "System Packages"
        sudo zypper update -y
    fi

    if command -v home-manager &>/dev/null; then
        print_step "Updating Nix packages & GUI apps"
        hm_target="tmojzes"
        if [ "$(uname)" = "Darwin" ]; then
            hm_target="tmojzes-mac"
        elif [ "$(uname -m)" = "x86_64" ]; then
            hm_target="tmojzes-x86_64-linux"
        fi
        nix flake update --flake "$HOME/.config/home-manager" &&
            home-manager switch --flake "$HOME/.config/home-manager#$hm_target"
    fi

    update_tool snap "Updating Snaps" sudo snap refresh
    update_tool go-global-update "Updating Go packages" go-global-update
    if command -v rustup &>/dev/null; then
        if ! rustup show active-toolchain &>/dev/null; then
            print_step "Installing Rust default toolchain"
            rustup default stable
        else
            print_step "Updating Rust toolchain"
            rustup update
        fi
    fi

    if cargo -V &>/dev/null; then
        print_step "Updating Rust crates"
        if ! command -v cargo-binstall &>/dev/null; then
            curl -fsSL https://raw.githubusercontent.com/cargo-bins/cargo-binstall/main/install-from-binstall-release.sh | bash
        fi
        if ! cargo install-update -V &>/dev/null; then
            if command -v cargo-binstall &>/dev/null; then
                cargo binstall -y cargo-update
            else
                cargo install cargo-update
            fi
        fi
        cargo install-update -a
    fi

    update_tool bob "Updating Bob" bob_update

    if command -v ibmcloud &>/dev/null && [ -d "$HOME/ibmcloud_homes" ]; then
        print_step "Updating IBM Cloud Plugins"

        ibmcloud plugin update --all -f

        for ic_home in "$HOME/ibmcloud_homes"/*; do
            if [ -d "$ic_home" ]; then
                echo "-> Profile: $(basename "$ic_home")"
                IBMCLOUD_HOME="$ic_home" ibmcloud plugin update --all -f
            fi
        done
    fi

    print_title "All updates complete!"
}
