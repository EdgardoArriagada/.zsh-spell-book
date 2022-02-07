alias cdzsb="cds ${ZSB_DIR}"
alias vzsb="cdzsb && nnvim"
alias galias="alias | rg"
alias his="nnvim ~/.zsh_history"
alias lis="tail -5 ~/.zsh_history | c -l ruby -"
alias hrc="nnvim ~/.zshrc"
alias env="nnvim ${ZSB_ZSHENV}"

hisIgnore cdzsb vzsb his hrc env lis
