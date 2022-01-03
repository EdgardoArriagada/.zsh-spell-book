cdg(){
  if [[ -n "$1" ]]; then
    eval "cd ${1}"; return 0
  fi

  local choosenPath=$(fd -t d | fzf)

  [[ -z "$choosenPath" ]] && return 0

  print -z "${0} ${choosenPath}"
}

_${zsb}.cdg() {
  local newCompletion=( $(fd -t d) )
  _describe 'command' newCompletion
}

compdef _${zsb}.cdg cdg
