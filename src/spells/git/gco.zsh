# I only use this for checking out files
gco() {
  git checkout "$@" && ${zsb}_gitStatus
}

complete -C "$ZSB_GIT_UNSTAGED_FILES" gco
