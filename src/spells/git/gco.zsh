# I only use this for checking out files
gco() {
  git checkout "$@" && ${zsb}.gitStatus
}

compdef "_${zsb}.gitUnrepeat 'unstaged'" gco

