_${zsb}.listFromCommand() {
  local usedcompletion=( "${words[@]:1:$CURRENT-2}" )
  local completionlist=( $(eval "$1") )

  local newcompletion=( ${completionlist:|usedcompletion} )
  _describe 'command' newcompletion
}

