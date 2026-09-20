# If not running interactively, don't do anything
[[ $- != *i* ]] && return

# Source global definitions
if [ -f /etc/zshrc ]; then
    . /etc/zshrc
fi

# User specific environment
if ! [[ "${PATH}" =~ "${HOME}/.local/bin:${HOME}/bin:" ]]; then
    PATH="${HOME}/.local/bin:${HOME}/bin:${PATH}"
fi
export PATH

# User specific aliases and functions (shared with bash via ~/.shrc.d)
if [ -d ~/.shrc.d ]; then
    for rc in ~/.shrc.d/*(N); do
        if [ -f "${rc}" ]; then
            . "${rc}"
        fi
    done
fi
unset rc
