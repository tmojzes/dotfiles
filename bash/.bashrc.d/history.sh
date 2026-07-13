export HISTFILESIZE=
export HISTSIZE=
export HISTTIMEFORMAT="[%F %T] "
# Change the file location, because certain Bash sessions truncate .bash_history file upon close.
# <https://superuser.com/questions/575479>
#   Bash history truncated to 500 lines on each login
#
export HISTFILE=~/.bash_eternal_history
# Force prompt to write history after every command.
# <https://superuser.com/questions/20900/bash-history-loss>
#   Bash history loss when using histappend
#
PROMPT_COMMAND="history -a; $PROMPT_COMMAND"
