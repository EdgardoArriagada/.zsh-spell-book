${zsb}.searchFileAndExecute() {
  local inputAlias="$1"
  local inputCommand="$2"
  local dataGetter="$3"
  shift 3

  if [[ -n "$@" ]]; then
    eval "${inputCommand} ${@}"; return $?
  fi

  local choosenFile=$(eval "$dataGetter" | fzf)
  [[ -z "$choosenFile" ]] && return 0

  print -z "${inputAlias} ${choosenFile}"
}

_${zsb}.searchFileAndExecute() {
  local newCompletion=( $(eval "$1") )
  _describe 'command' newCompletion
}

compdef _${zsb}.searchFileAndExecute ${zsb}.searchFileAndExecute

vg_dataGetter='fd -t f'
vg() { ${zsb}.searchFileAndExecute "$0" v "$vg_dataGetter" "$@" }
compdef "_${zsb}.searchFileAndExecute '${vg_dataGetter}'" vg

cg_dataGetter='fd -t f'
cg() { ${zsb}.searchFileAndExecute "$0" c "$cg_dataGetter" "$@" }
compdef "_${zsb}.searchFileAndExecute '${cg_dataGetter}'" cg

cdg_dataGetter='fd -t d'
cdg() { ${zsb}.searchFileAndExecute "$0" cd "$cdg_dataGetter" "$@" }
compdef "_${zsb}.searchFileAndExecute '${cdg_dataGetter}'" cdg

cdd_dataGetter='dirs -p'
cdd() { ${zsb}.searchFileAndExecute "$0" cd "$cdd_dataGetter" "$@" }
compdef "_${zsb}.searchFileAndExecute '${cdd_dataGetter}'" cdd

