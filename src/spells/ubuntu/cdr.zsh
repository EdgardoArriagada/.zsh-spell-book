cdr() {
  if [[ -n "$1" ]]; then
    eval "cd ${1}"
    return $?
  fi

  local choosenDir=$(dirs -p | fzf)
  [[ -z "$choosenDir" ]] && return 0
  eval "cd ${choosenDir}"
}

_${zsb}.cdr() {
  local newCompletion=( $(dirs -p) )

  _${zsb}.verticalComp "newCompletion"
}

compdef _${zsb}.cdr cdr
