${zsb}.searchFileAndOpen() {
  [[ -z "$2" ]] && ${zsb}.throw "You have to pass an argument"
  eval "${@}" && return $?
}

_${zsb}.searchFileAndOpen() {
  local newCompletion=( $(fd -t f) )
  _describe 'command' newCompletion
}

compdef _${zsb}.searchFileAndOpen ${zsb}.searchFileAndOpen

alias cg="${zsb}.searchFileAndOpen c"
alias vg="${zsb}.searchFileAndOpen v"
