gsl() {
  local gslEmptyMsg="Stash list is empty."
  local -r stashList="`git stash list`"

  [[ -z "$stashList" ]] && ${zsb}.cancel "$gslEmptyMsg"

  if [[ "$1" = "--clear" ]]; then
    ${zsb}.warning "You are about to clear the `hl 'git stash list'`"
    ${zsb}.confirmMenu.withPrompt

    git stash clear && ${zsb}.success "$gslEmptyMsg"
    return 0
  fi

  echo "$stashList" | c -p -l ruby -
}

hisIgnore gsl 'gsl --clear'

compdef "_${zsb}.singleComp '--clear'" gsl

alias GSL="toggleCapsLock && gsl"
