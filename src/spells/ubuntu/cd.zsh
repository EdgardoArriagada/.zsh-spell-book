cd() {
  builtin cd $1 && ls

  [[ -f .nvmrc ]] && eval 'nvm use'
  return 0
}

alias CD='toggleCapsLock && cd'

alias dc='cd'
alias DC='CD'
