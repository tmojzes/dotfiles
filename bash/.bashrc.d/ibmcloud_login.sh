ibmcloud_login() {
    if ! command -v ibmcloud &>/dev/null; then
        echo "Error: ibmcloud cli not found"
        return 1
    fi

    print_title() { echo -e "\n\033[1;34m--- $1 ---\033[0m"; }

    print_title "Logging into IBM Cloud Envs"

    local base_dir="$HOME/ibmcloud_homes"
    local region="us-south"

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

        # 2. Loop through the API keys inside this endpoint folder
        for api_key_file in "$endpoint_dir"/*; do
            # Skip if it's not a regular file
            [ -f "$api_key_file" ] || continue

            local profile=$(basename "$api_key_file")
            local home_dir="$base_dir/$profile"

            echo "-> Logging into profile: $profile ($endpoint)"

            # Ensure the profile's home directory exists
            mkdir -p "$home_dir"

            # Execute the login
            IBMCLOUD_HOME="$home_dir" ibmcloud login -q --apikey @"$api_key_file" -a "$endpoint" -r "$region"
        done
    done

    print_title "Completed Logins"
}
