gb() {
  case $# in
    0) <<< '' ; git branch; <<< '' ;;
    *) git branch "$@" ;;
  esac
}

compdef "_${zsb}.nonRepeatedListC ${zsb}.gitBranches" gb

alias GB="toggleCapsLock && gb"

