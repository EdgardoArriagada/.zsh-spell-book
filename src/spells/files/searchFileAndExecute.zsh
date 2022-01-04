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
vg_command='v'
vg() ${zsb}.searchFileAndExecute "$0" "$vg_command" "$vg_dataGetter" "$@"
compdef "_${zsb}.searchFileAndExecute '${vg_dataGetter}'" vg

cg_dataGetter='fd -t f'
cg_command='c'
cg() ${zsb}.searchFileAndExecute "$0" "$cg_command" "$cg_dataGetter" "$@"
compdef "_${zsb}.searchFileAndExecute '${cg_dataGetter}'" cg

cdg_dataGetter='fd -t d'
cdg_command='cd'
cdg() ${zsb}.searchFileAndExecute "$0" "$cdg_command" "$cdg_dataGetter" "$@"
compdef "_${zsb}.searchFileAndExecute '${cdg_dataGetter}'" cdg

cdd_dataGetter='dirs -p'
cdd_command='cd'
cdd() ${zsb}.searchFileAndExecute "$0" "$cdd_command" "$cdd_dataGetter" "$@"
compdef "_${zsb}.searchFileAndExecute '${cdd_dataGetter}'" cdd

