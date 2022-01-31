alias cdzsb="cds ${ZSB_DIR}"
alias vzsb="cdzsb && nnvim"
alias galias="alias | rg"
alias his="nnvim ~/.zsh_history"
alias hrc="nnvim ~/.zshrc"
alias env="nnvim ${ZSB_ZSHENV}"

hisIgnore cdzsb vzsb his hrc env
