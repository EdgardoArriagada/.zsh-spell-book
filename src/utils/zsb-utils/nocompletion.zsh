_${zsb}.acceptNoArgs() {
  [ "$CURRENT" -gt "1" ] && return 0
}


_${zsb}.nocompletion() {
  compdef _${zsb}.acceptNoArgs $1
}

