alias cdzsb="cds $ZSB_DIR"
alias vzsb="builtin cd $ZSB_DIR && nvim && ${zsb}.gitStatus"
alias galias="alias | rg"
alias his="nvim ~/.zsh_history"
alias lis="tail -5 ~/.zsh_history | c -p -l ruby -"
alias hrc="nvim ~/.zshrc"
alias env="nvim $ZSB_ZSHENV"

hisIgnore cdzsb vzsb his hrc env lis
