gd() {
  if [ -z "$1" ]; then
    git diff && ${zsb}_gitStatus
    return 0
  fi

  git diff "$@" && ${zsb}_gitStatus
  return 0
}

complete -C "$ZSB_GIT_UNSTAGED_FILES" gd
