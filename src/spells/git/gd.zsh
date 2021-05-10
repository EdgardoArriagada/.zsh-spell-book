gd() {
  if [ -z "$1" ]; then
    git diff && ${zsb}.gitStatus
    return 0
  fi

  git diff "$@" && ${zsb}.gitStatus
  return 0
}

compdef "_${zsb}.gitUnrepeat 'unstaged'" gd

