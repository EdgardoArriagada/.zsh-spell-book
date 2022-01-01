${zsb}.searchFileAndOpen() {
  [[ -n "$2" ]] && eval "${@}" && return $?

  local choosenDir=$(fd -t f | fzf)
  [[ -z "$choosenDir" ]] && return 0
  eval "${1} ${choosenDir}"
}

_${zsb}.searchFileAndOpen() {
  local newCompletion=( $(fd -t f) )

  _describe 'command' newCompletion
}

compdef _${zsb}.searchFileAndOpen ${zsb}.searchFileAndOpen

alias cg="${zsb}.searchFileAndOpen c"
alias vg="${zsb}.searchFileAndOpen v"
