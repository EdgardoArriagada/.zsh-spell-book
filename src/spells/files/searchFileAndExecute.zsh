${zsb}.searchFileAndExecute() {
  [[ -n "$2" ]] && eval "${@}" && return $?

  local choosenFile=$(fd -t f | fzf)

  [[ -z "$choosenFile" ]] && return 0

  print -z "${1}g ${choosenFile}"
}

_${zsb}.searchFileAndExecute() {
  local newCompletion=( $(fd -t f) )
  _describe 'command' newCompletion
}

compdef _${zsb}.searchFileAndExecute ${zsb}.searchFileAndExecute

alias cg="${zsb}.searchFileAndExecute c"
alias vg="${zsb}.searchFileAndExecute v"
