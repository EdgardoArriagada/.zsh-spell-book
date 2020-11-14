gco() {
  git checkout "$@"
}

complete -C "$ZSB_GIT_UNSTAGED_FILES" gco
