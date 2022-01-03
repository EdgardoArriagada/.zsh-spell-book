cdr(){
  if [[ -n "$1" ]]; then
    eval "cd ${1}"; return 0
  fi

  local choosenPath=$(dirs -p | fzf)

  [[ -z "$choosenPath" ]] && return 0

  print -z "${0} ${choosenPath}"
}

_${zsb}.cdr() {
  local newCompletion=( $(dirs -p) )
  _describe 'command' newCompletion
}

compdef _${zsb}.cdr cdr
