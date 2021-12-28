cd() {
  builtin cd $1 && ls

  [[ -f .nvmrc ]] && eval 'nvm use'
}

alias CD='toggleCapsLock && cd'

alias dc='cd'
alias DC='CD'
