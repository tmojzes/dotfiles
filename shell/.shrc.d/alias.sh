# Set aliases
alias l='ls -l --color=auto'
alias ll='ls -lah --color=auto'
alias grep='grep --color=auto'
alias vim="nvim"
alias gemini='npx @google/gemini-cli@latest'
alias k="kubectl"

# Search for a file by name pattern (function instead of alias so the argument works).
rgf() { rg --files | rg "$1"; }

# Using functions instead of alias makes completion possible.
if command -v ibmcloud &>/dev/null; then
    ic() { IBMCLOUD_HOME="" ibmcloud "$@"; }
    ich() { IBMCLOUD_HOME="$HOME/ibmcloud_homes/hl_dev" ibmcloud "$@"; }
    icp() { IBMCLOUD_HOME="$HOME/ibmcloud_homes/personal" ibmcloud "$@"; }
    icsm() { IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_mesh" ibmcloud "$@"; }
    icst() { IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_test" ibmcloud "$@"; }
    icstl() { IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_test_location" ibmcloud "$@"; }
    icss() { IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_stage" ibmcloud "$@"; }
    icad() { IBMCLOUD_HOME="$HOME/ibmcloud_homes/argonauts_dev" ibmcloud "$@"; }
    ihld() { IBMCLOUD_HOME="$HOME/ibmcloud_homes/hybrid_link_dev" ibmcloud "$@"; }
    ihlp() { IBMCLOUD_HOME="$HOME/ibmcloud_homes/hybrid_link_prod" ibmcloud "$@"; }

    if [ -n "$ZSH_VERSION" ]; then
        # 2>/dev/null: compdef errors on hosts where the ibmcloud zsh autocomplete service is absent.
        compdef ic=ibmcloud icp=ibmcloud ich=ibmcloud icsm=ibmcloud icst=ibmcloud \
            icstl=ibmcloud icss=ibmcloud icad=ibmcloud ihld=ibmcloud ihlp=ibmcloud 2>/dev/null
    else
        complete -o default -F _bash_autocomplete ic
        complete -o default -F _bash_autocomplete ich
        complete -o default -F _bash_autocomplete icp
        complete -o default -F _bash_autocomplete icsm
        complete -o default -F _bash_autocomplete icst
        complete -o default -F _bash_autocomplete icstl
        complete -o default -F _bash_autocomplete icss
        complete -o default -F _bash_autocomplete icad
        complete -o default -F _bash_autocomplete ihld
        complete -o default -F _bash_autocomplete ihlp
    fi
fi
# alias sopse="sops --encrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"
# alias sopsd="sops --decrypt --age $(cat $SOPS_AGE_KEY_FILE | grep -oP 'public key: \K(.*)') -i"
