# edge case: tryng to add new file named 'new' using 'add new'
# will trigger 'git add .' instead of 'git add new'

ga() {
  if [ "$1" = "new" ] || [ "$1" = "." ]; then
    if [ -z $2 ]; then
      git add . && gitStatus
      return 0
    fi
    shift 1
    git add "$@" && gitStatus
    return 0
  fi

  git add -p "$@" && gitStatus
  return 0
}

# Complete with unstaged/untracked files (from current dir)
complete -C "git status --short | grep -E '(^ [MARCD\?])|(^[MARCD\?]{2})' | sed s/^...//" ga
