# I only use this to move staged files back to unstaged state
grh() {
  git reset -q "$@" && ${zsb}_gitStatus
}


complete -C "$ZSB_GIT_STAGED_FILES" grh
