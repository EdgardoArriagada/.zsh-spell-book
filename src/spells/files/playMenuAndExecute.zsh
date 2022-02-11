${zsb}.playMenuAndExecute() {
  local inputAlias="$1"
  local inputCommand="$2"
  local dataGetter="$3"
  shift 3

  if [[ -a "$1" ]]; then
    eval "${inputCommand} ${@}"; return $?
  fi

  local choosenFile=$(eval "$dataGetter" | fzf --query "$1")
  [[ -z "$choosenFile" ]] && return 0

  print -z "${inputAlias} ${choosenFile}"
}

_${zsb}.playMenuAndExecute() {
  local newCompletion=( $(eval "$1") )
  _describe 'command' newCompletion
}

compdef _${zsb}.playMenuAndExecute ${zsb}.playMenuAndExecute

__vg_dataGetter='fd -t f'
__vg_command='v'
vg() ${zsb}.playMenuAndExecute "$0" "$__vg_command" "$__vg_dataGetter" "$@"
compdef "_${zsb}.playMenuAndExecute '${__vg_dataGetter}'" vg
hisIgnore vg

__vgh_dataGetter='fd -t f --hidden --no-ignore'
__vgh_command='v'
vgh() ${zsb}.playMenuAndExecute "$0" "$__vgh_command" "$__vgh_dataGetter" "$@"
_${zsb}.nocompletion cgh
hisIgnore vgh

__cg_dataGetter='fd -t f'
__cg_command='c'
cg() ${zsb}.playMenuAndExecute "$0" "$__cg_command" "$__cg_dataGetter" "$@"
compdef "_${zsb}.playMenuAndExecute '${__cg_dataGetter}'" cg
hisIgnore cg

__cgh_dataGetter='fd -t f --hidden --no-ignore'
__cgh_command='c'
cgh() ${zsb}.playMenuAndExecute "$0" "$__cgh_command" "$__cgh_dataGetter" "$@"
_${zsb}.nocompletion cgh
hisIgnore cgh

__cdg_dataGetter='fd -t d'
__cdg_command='cd'
cdg() ${zsb}.playMenuAndExecute "$0" "$__cdg_command" "$__cdg_dataGetter" "$@"
compdef "_${zsb}.playMenuAndExecute '${__cdg_dataGetter}'" cdg
hisIgnore cdg

__cdgh_dataGetter='fd -t d --hidden --no-ignore'
__cdgh_command='cd'
cdgh() ${zsb}.playMenuAndExecute "$0" "$__cdgh_command" "$__cdgh_dataGetter" "$@"
_${zsb}.nocompletion cdgh
hisIgnore cdgh

__cdd_dataGetter='dirs -p'
__cdd_command='cd'
cdd() ${zsb}.playMenuAndExecute "$0" "$__cdd_command" "$__cdd_dataGetter" "$@"
compdef "_${zsb}.playMenuAndExecute '${__cdd_dataGetter}'" cdd
hisIgnore cdd

__cdp_dataGetter='fd'
__cdp_command='cd'
cdp() {
  if ! [[ -a "$1" ]]; then
    ${zsb}.playMenuAndExecute "$0" "$__cdp_command" "$__cdp_dataGetter" "$@"
    return 0
  fi

  if [[ -d "$1" ]]; then
    eval "cd ${1}"; return 0
  fi

  if [[ -f "$1" ]]; then
    eval "cd $(dirname ${1})"; return 0
  fi
}
compdef "_${zsb}.playMenuAndExecute '${__cdp_dataGetter}'" cdp
hisIgnore cdp

