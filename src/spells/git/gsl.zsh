gsl() {
  ${zsb}.validateStashList

  if [[ "$1" = "--clear" ]]; then
    ${zsb}.confirmMenu.warning "You are about to clear the `hl 'git stash list'`"

    git stash clear && ${zsb}.success 'Stash list is empty'
    return 0
  fi

  git stash list | c -p -l ruby -
}

hisIgnore gsl 'gsl --clear'

compdef "_${zsb}.singleComp '--clear'" gsl

alias GSL="toggleCapsLock && gsl"
