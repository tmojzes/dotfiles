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

    bob_update() {
        version=$(curl -s https://s3.us-south.cloud-object-storage.appdomain.cloud/bobshell/bobshell-version.txt)
        dl_url="https://s3.us-south.cloud-object-storage.appdomain.cloud/bobshell/bobshell-${version}.tgz"

        npm install --reg=https://registry.npmjs.org/ -g "${dl_url}"
    }

    update_tool snap "Updating Snaps" sudo snap refresh
    update_tool flatpak "Updating Flatpaks" flatpak update -y
    update_tool go-global-update "Updating Go packages" go-global-update
    update_tool rustup "Updating Rust toolchain" rustup update
    update_tool bob "Updating Bob" bob_update

    if command -v home-manager &>/dev/null; then
        print_step "Updating Nix packages"
        hm_target="tmojzes"
        if [ "$(uname)" = "Darwin" ]; then
            hm_target="tmojzes-mac"
        elif [ "$(uname -m)" = "x86_64" ]; then
            hm_target="tmojzes-x86_64-linux"
        fi
        nix flake update --flake "$HOME/.config/home-manager" &&
            home-manager switch --flake "$HOME/.config/home-manager#$hm_target"
    fi

    if command -v cargo &>/dev/null; then
        print_step "Updating Rust crates"
        if ! cargo install-update -V &>/dev/null; then
            cargo install cargo-update
        fi
        cargo install-update -a
    fi

    if command -v ibmcloud &>/dev/null && [ -d "$HOME/ibmcloud_homes" ]; then
        print_step "Updating IBM Cloud Plugins"

        ibmcloud update plugin update --all -f

        for ic_home in "$HOME/ibmcloud_homes"/*; do
            if [ -d "$ic_home" ]; then
                echo "-> Profile: $(basename "$ic_home")"
                IBMCLOUD_HOME="$ic_home" ibmcloud plugin update --all -f
            fi
        done
    fi

    print_title "All updates complete!"
}
