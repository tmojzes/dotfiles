ibmcloud_login() {
    # Check the alias.sh for ibmcloud command aliases.
    echo "--- Logging into IBM Cloud Envs  ---"

    ich login -q --apikey @~/ibmcloud_homes/api_keys/hl_dev -a https://test.cloud.ibm.com -r "us-south"
    icp login -q --apikey @~/ibmcloud_homes/api_keys/personal -a https://cloud.ibm.com -r "us-south"

    echo "--- Completed logins  ---"

}
