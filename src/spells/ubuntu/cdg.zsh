cdg() eval "cd ${1}"

_${zsb}.cdg() {
  local newCompletion=( $(fd -t d) )
  _describe 'command' newCompletion
}

compdef _${zsb}.cdg cdg
