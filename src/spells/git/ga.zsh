# edge case: tryng to add new file named 'new' using 'add new'
# will trigger 'git add .' instead of 'git add new'

ga() {
  if [ "$1" = "new" ] || [ "$1" = "." ]; then
    if [ -z "$2" ]; then
      git add . && ${zsb}_gitStatus
      return 0
    fi
    shift 1
    git add "$@" && ${zsb}_gitStatus
    return 0
  fi

  git add -p "$@" && ${zsb}_gitStatus
  return 0
}

complete -C "$ZSB_GIT_UNSTAGED_AND_NEW_FILES" ga
