grh() {
  git reset -q "$@" && ${zsb}_gitStatus
}


complete -C "$ZSB_GIT_STAGED_FILES" grh
