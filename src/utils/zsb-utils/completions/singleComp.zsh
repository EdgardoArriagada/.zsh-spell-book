_${zsb}.singleComp() {
  [[ "$CURRENT" -gt 2 ]] && return 0
  local -r comp=( "$@" )
  _describe 'command' comp
}

