# I only use this to move staged files back to unstaged state
grh() {
  git reset -q "$@" && ${zsb}.gitStatus
}

compdef "_${zsb}.gitUnrepeat 'staged'" grh

