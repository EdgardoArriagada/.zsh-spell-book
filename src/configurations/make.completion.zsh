_${zsb}.make() {
  [[ "$CURRENT" -gt 2 ]] && return 0
  [[ ! -f ./Makefile ]] && return 0

  local -r parsedMakefile="`
      grep --only-matching '^[a-zA-Z-]*:' Makefile |\
      sd ':' ''
  `"

  local -r compList=( "${(@f)parsedMakefile}" )

  _describe 'command' compList
}

compdef _${zsb}.make make

