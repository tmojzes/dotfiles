# Set aliases
alias l='ls -l --color=auto'
alias ll='ls -lah --color=auto'
alias grep='grep --color=auto'
alias vim="nvim"
alias gemini='npx @google/gemini-cli@latest'

# ibmcloud aliases
alias ic='ibmcloud'
alias ich='IBMCLOUD_HOME=~/ibmcloud_homes/hl_dev ibmcloud'
alias icp='IBMCLOUD_HOME=~/ibmcloud_homes/personal ibmcloud'

# alias sopse="sops --encrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"
# alias sopsd="sops --decrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"
