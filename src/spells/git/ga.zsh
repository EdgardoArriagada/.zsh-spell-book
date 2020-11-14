# edge case: tryng to add new file named 'new' using 'add new'
# will trigger 'git add .' instead of 'git add new'

ga() {
  if [ "$1" = "." ]; then
    git add . && ${zsb}_gitStatus
    return 0
  fi

  if [ "$1" = "new" ]; then
    if [ -z "$2" ]; then
      git add . && ${zsb}_gitStatus
      return 0
    fi
    shift 1 # remove 'new' flag
    git add "$@" && ${zsb}_gitStatus
    return 0
  fi

  git add -p "$@" && ${zsb}_gitStatus
  return 0
}

_${zsb}_ga() {
  if [ "${COMP_WORDS[1]}" = "." ]; then
    return 0
  fi

  if [ "${COMP_WORDS[1]}" = "new" ]; then
    COMPREPLY=( $(compgen -C "$ZSB_GIT_UNSTAGED_AND_UNTRACKED_FILES") )
    return 0
  fi

  COMPREPLY=( $(compgen -C "$ZSB_GIT_UNSTAGED_FILES") )
}

complete -F _${zsb}_ga ga
