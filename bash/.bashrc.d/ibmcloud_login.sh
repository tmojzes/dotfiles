ibmcloud_login() {
    if ! command -v ibmcloud &>/dev/null; then
        echo "Error: ibmcloud cli not found"
        return 1
    fi

    print_title() { echo -e "\n\033[1;34m--- $1 ---\033[0m"; }

    print_title "Logging into IBM Cloud Envs"

    local base_dir="$HOME/ibmcloud_homes"
    local region="us-south"
    local failed=0

    if [ ! -d "$base_dir/api_keys" ]; then
        echo "Error: Directory $base_dir/api_keys not found."
        return 1
    fi

    # 1. Loop through Endpoint directories (e.g., cloud.ibm.com)
    for endpoint_dir in "$base_dir/api_keys"/*; do
        # Skip if it's not a directory
        [ -d "$endpoint_dir" ] || continue

        # The folder name IS the endpoint!
        local endpoint="https://$(basename "$endpoint_dir")"
        local containers_host

        case "$endpoint" in
        https://cloud.ibm.com)
            containers_host="https://containers.cloud.ibm.com"
            ;;
        https://test.cloud.ibm.com)
            containers_host="https://containers.test.cloud.ibm.com"
            ;;
        *)
            containers_host=""
            ;;
        esac

        # 2. Loop through the API keys inside this endpoint folder
        for api_key_file in "$endpoint_dir"/*; do
            # Skip if it's not a regular file
            [ -f "$api_key_file" ] || continue

            local profile=$(basename "$api_key_file")
            local home_dir="$base_dir/$profile"

            echo "-> Logging into profile: $profile ($endpoint)"

            # Ensure the profile's home directory exists
            mkdir -p "$home_dir"

            # Execute the login.
            if ! IBMCLOUD_HOME="$home_dir" ibmcloud login -q --apikey @"$api_key_file" -a "$endpoint" -r "$region"; then
                echo "Error: Login failed for profile: $profile" >&2
                failed=1
                continue
            fi

            if [ -n "$containers_host" ]; then
                echo "-> Initializing Container Service endpoint: $containers_host"
                if ! IBMCLOUD_HOME="$home_dir" ibmcloud ks init --host "$containers_host" -q; then
                    echo "Error: Container Service initialization failed for profile: $profile" >&2
                    failed=1
                fi
            fi
        done
    done

    print_title "Completed Logins"
    return "$failed"
}
