gb() {
  if [[ -z "$1" ]]; then
    echo "\n"; git branch; echo "\n"; return 0
  fi

  git branch "$@"
}

_gb() { git branch | sd "^.." ""; }

compdef "_${zsb}.nonRepeatedListC _gb" gb

alias GB="toggleCapsLock && gb"

