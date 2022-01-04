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

vgh_dataGetter='fd -t f --hidden --no-ignore'
vgh_command='v'
vgh() ${zsb}.searchFileAndExecute "$0" "$vgh_command" "$vgh_dataGetter" "$@"
_${zsb}.nocompletion cgh

cg_dataGetter='fd -t f'
cg_command='c'
cg() ${zsb}.searchFileAndExecute "$0" "$cg_command" "$cg_dataGetter" "$@"
compdef "_${zsb}.searchFileAndExecute '${cg_dataGetter}'" cg

cgh_dataGetter='fd -t f --hidden --no-ignore'
cgh_command='c'
cgh() ${zsb}.searchFileAndExecute "$0" "$cgh_command" "$cgh_dataGetter" "$@"
_${zsb}.nocompletion cgh

cdg_dataGetter='fd -t d'
cdg_command='cd'
cdg() ${zsb}.searchFileAndExecute "$0" "$cdg_command" "$cdg_dataGetter" "$@"
compdef "_${zsb}.searchFileAndExecute '${cdg_dataGetter}'" cdg

cdgh_dataGetter='fd -t d --hidden --no-ignore'
cdgh_command='cd'
cdgh() ${zsb}.searchFileAndExecute "$0" "$cdgh_command" "$cdgh_dataGetter" "$@"
_${zsb}.nocompletion cdgh

cdd_dataGetter='dirs -p'
cdd_command='cd'
cdd() ${zsb}.searchFileAndExecute "$0" "$cdd_command" "$cdd_dataGetter" "$@"
compdef "_${zsb}.searchFileAndExecute '${cdd_dataGetter}'" cdd

