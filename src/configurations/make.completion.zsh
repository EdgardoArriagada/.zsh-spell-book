_${zsb}.make() {
  [[ "$CURRENT" -gt 2 ]] && return 0
  [[ ! -f ./Makefile ]] && return 0

  local -r parsedMakefile="$(
    grep -C 1 '^[a-zA-Z]' Makefile |\
    sed -E ':a ; $!N ; s/\n\s+// ; ta ; P ; D' |\
    sed -E ':a;$!N;s/\n--/ \(\.\.\.\)/;ta;P;D'
  )"

  local -r compList=( "${(@f)parsedMakefile}" )

  _describe 'command' compList
}

compdef _${zsb}.make make

