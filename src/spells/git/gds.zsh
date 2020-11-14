gds() {
  if [ -z "$1" ]; then
    git diff --staged && ${zsb}_gitStatus
    return 0
  fi

  git diff --staged "$@" && ${zsb}_gitStatus
  return 0
}

complete -C "$ZSB_GIT_STAGED_FILES" gds
