cdg() {
  if [[ -n "$1" ]]; then
    eval "cd ${1}"
    return $?
  fi

  local choosenDir=$(fd -t d | fzf)
  [[ -z "$choosenDir" ]] && return 0
  eval "cd ${choosenDir}"
}

_${zsb}.cdg() {
  local newCompletion=( $(fd -t d) )

  _${zsb}.verticalComp "newCompletion"
}

compdef _${zsb}.cdg cdg
