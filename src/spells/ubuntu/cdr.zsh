cdr() eval "cd ${1}"

_${zsb}.cdr() {
  local newCompletion=( $(dirs -p) )
  _describe 'command' newCompletion
}

compdef _${zsb}.cdr cdr
