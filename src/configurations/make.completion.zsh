_${zsb}.make() {
  (( $CURRENT > 2 )) && return 0
  [[ ! -f ./Makefile ]] && return 0

  local parsedMakefile="`
      grep --only-matching '^[a-zA-Z-]*:' Makefile |\
      sd ':' ''
  `"

  local compList=( "${(@f)parsedMakefile}" )

  _describe 'command' compList
}

compdef _${zsb}.make make

