# Set aliases
alias l='ls -l --color=auto'
alias ll='ls -lah --color=auto'
alias grep='grep --color=auto'
alias vim="nvim"
alias gemini='npx @google/gemini-cli@latest'

# Using functions instead of alias makes bash completion possible.
if command -v ibmcloud &>/dev/null; then
    ic() {
        IBMCLOUD_HOME="" ibmcloud "$@"
    }
    complete -o default -F _bash_autocomplete ic

    ich() {
        IBMCLOUD_HOME="$HOME/ibmcloud_homes/hl_dev" ibmcloud "$@"
    }
    complete -o default -F _bash_autocomplete ich

    icp() {
        IBMCLOUD_HOME="$HOME/ibmcloud_homes/personal" ibmcloud "$@"
    }
    complete -o default -F _bash_autocomplete icp

    icsm() {
        IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_mesh" ibmcloud "$@"
    }
    complete -o default -F _bash_autocomplete icsm

    icst() {
        IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_test" ibmcloud "$@"
    }
    complete -o default -F _bash_autocomplete icst

    icstl() {
        IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_test_location" ibmcloud "$@"
    }
    complete -o default -F _bash_autocomplete icstl

    icss() {
        IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_stage" ibmcloud "$@"
    }
    complete -o default -F _bash_autocomplete icss

    icad() {
        IBMCLOUD_HOME="$HOME/ibmcloud_homes/argonauts_dev" ibmcloud "$@"
    }
    complete -o default -F _bash_autocomplete icad

fi
# alias sopse="sops --encrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"
# alias sopsd="sops --decrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"
