gb() {
  case $# in
    0) <<< '' ; git branch; <<< ''  ;;
    *) git branch "$@" ;;
  esac
}

_gb() { git branch | sd "^.." ""; }

compdef "_${zsb}.nonRepeatedListC _gb" gb

alias GB="toggleCapsLock && gb"

