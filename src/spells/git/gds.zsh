gds() {
  if [ -z "$1" ]; then
    git diff --staged && gitStatus
    return 0
  fi

  git diff --staged "$@" && gitStatus
  return 0
}

# Complete only with staged files
complete -C "git status --short | grep '^[MARCD]' | sed s/^...//" gds
