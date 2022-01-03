${zsb}.searchFileAndOpen() {
  [[ -n "$2" ]] && eval "${@}" && return $?

  local choosenFile=$(fd -t f | fzf)

  [[ -z "$choosenFile" ]] && return 0

  print -z "${1}g ${choosenFile}"
}

_${zsb}.searchFileAndOpen() {
  local newCompletion=( $(fd -t f) )
  _describe 'command' newCompletion
}

compdef _${zsb}.searchFileAndOpen ${zsb}.searchFileAndOpen

alias cg="${zsb}.searchFileAndOpen c"
alias vg="${zsb}.searchFileAndOpen v"
