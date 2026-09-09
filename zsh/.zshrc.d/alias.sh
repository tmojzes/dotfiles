# Set aliases
alias l='ls -l --color=auto'
alias ll='ls -lah --color=auto'
alias grep='grep --color=auto'
alias vim="nvim"
alias gemini='npx @google/gemini-cli@latest'
alias k="kubectl"
alias rgf="rg --files | rg "$1""

# Using functions instead of alias makes zsh completion possible.
# 'function' keyword prevents alias expansion during parsing.
if (("$+commands[ibmcloud]")); then
    function ic()    { IBMCLOUD_HOME="" ibmcloud "$@" }
    function ich()   { IBMCLOUD_HOME="$HOME/ibmcloud_homes/hl_dev" ibmcloud "$@" }
    function icp()   { IBMCLOUD_HOME="$HOME/ibmcloud_homes/personal" ibmcloud "$@" }
    function icsm()  { IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_mesh" ibmcloud "$@" }
    function icst()  { IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_test" ibmcloud "$@" }
    function icstl() { IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_test_location" ibmcloud "$@" }
    function icss()  { IBMCLOUD_HOME="$HOME/ibmcloud_homes/sat_stage" ibmcloud "$@" }
    function icad()  { IBMCLOUD_HOME="$HOME/ibmcloud_homes/argonauts_dev" ibmcloud "$@" }
    function ihld()  { IBMCLOUD_HOME="$HOME/ibmcloud_homes/hybrid_link_dev" ibmcloud "$@" }
    function ihlp()  { IBMCLOUD_HOME="$HOME/ibmcloud_homes/hybrid_link_prod" ibmcloud "$@" }

    compdef ic=ibmcloud icp=ibmcloud ich=ibmcloud icsm=ibmcloud icst=ibmcloud \
            icstl=ibmcloud icss=ibmcloud icad=ibmcloud ihld=ibmcloud ihlp=ibmcloud
fi
