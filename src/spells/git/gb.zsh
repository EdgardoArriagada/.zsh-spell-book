gb() {
  if [[ -z "$1" ]]; then
    echo "\n"; git branch; echo "\n"
    return 0
  fi

  git branch "$@"
}

compdef _git gb=git-branch

alias GB="toggleCapsLock && gb"

