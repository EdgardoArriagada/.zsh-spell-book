gds() {
  if [ -z "$1" ]; then
    git diff --staged && ${zsb}.gitStatus
    return 0
  fi

  git diff --staged "$@" && ${zsb}.gitStatus
  return 0
}

compdef "_${zsb}.gitUnrepeat 'staged'" gds

