cd() {
  builtin cd $1 && ls

  [[ -f .nvmrc ]] && ${zsb}.useCachedNodeVersion
  return 0
}

hisIgnore cd

alias CD='toggleCapsLock && cd'

alias dc='cd'
alias DC='CD'
