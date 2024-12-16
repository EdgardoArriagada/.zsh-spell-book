${zsb}.playMenuAndExecute() {
  local inputAlias=$1
  local inputCommand=$2
  local dataGetter=$3
  shift 3

  if [[ -a "$1" ]]; then
    eval "$inputCommand $@"; return $?
  fi

  local choosenFile=$(eval $dataGetter | fzf --query $1)
  [[ -z "$choosenFile" ]] && return 0

  print -z "$inputAlias $choosenFile"
}

_${zsb}.playMenuAndExecute() {
  local usedCompletion=( ${words[@]:1} )
  local cmd="${@} ${usedCompletion}"
  local newCompletion=( `eval ${cmd}` )
  _describe 'command' newCompletion
}

__vg_dataGetter='fd -t f --no-ignore'
__vg_command='v'
__vg() ${zsb}.playMenuAndExecute $0 $__vg_command $__vg_dataGetter $@
alias vg='noglob __vg'
compdef "_${zsb}.playMenuAndExecute $__vg_dataGetter" __vg
hisIgnore vg

__vgh_dataGetter='fd -t f --hidden'
__vgh_command='v'
__vgh() ${zsb}.playMenuAndExecute $0 $__vgh_command $__vgh_dataGetter $@
alias vgh='noglob __vgh'
compdef "_${zsb}.playMenuAndExecute $__vgh_dataGetter" __vgh
hisIgnore vgh

__tigg_dataGetter='fd -t f --no-ignore'
__tigg_command='tig'
__tigg() ${zsb}.playMenuAndExecute $0 $__tigg_command $__tigg_dataGetter $@
alias tigg='noglob __tigg'
compdef "_${zsb}.playMenuAndExecute $__tigg_dataGetter" __tigg
hisIgnore tigg

__tiggh_dataGetter='fd -t f --hidden'
__tiggh_command='tig'
__tiggh() ${zsb}.playMenuAndExecute $0 $__tiggh_command $__tiggh_dataGetter $@
alias tiggh='noglob __tiggh'
compdef "_${zsb}.playMenuAndExecute $__tiggh_dataGetter" __tiggh
hisIgnore tiggh

__cg_dataGetter='fd -t f --no-ignore'
__cg_command='c'
__cg() ${zsb}.playMenuAndExecute $0 $__cg_command $__cg_dataGetter $@
alias cg='noglob __cg'
compdef "_${zsb}.playMenuAndExecute $__cg_dataGetter" __cg
hisIgnore cg

__cgh_dataGetter='fd -t f --hidden'
__cgh_command='c'
__cgh() ${zsb}.playMenuAndExecute $0 $__cgh_command $__cgh_dataGetter $@
alias cgh='noglob __cgh'
hisIgnore cgh

__cdg_dataGetter='fd -t d --no-ignore'
__cdg_command='cd'
__cdg() ${zsb}.playMenuAndExecute $0 $__cdg_command $__cdg_dataGetter $@
alias cdg='noglob __cdg'
compdef "_${zsb}.playMenuAndExecute $__cdg_dataGetter" __cdg
hisIgnore cdg

__cdgh_dataGetter='fd -t d --hidden'
__cdgh_command='cd'
__cdgh() ${zsb}.playMenuAndExecute $0 $__cdgh_command $__cdgh_dataGetter $@
alias cdgh='noglob __cdgh'
compdef "_${zsb}.playMenuAndExecute $__cdgh_dataGetter" __cdgh
hisIgnore cdgh

__cdp_dataGetter='fd --no-ignore'
__cdp_command='cd'
__cdp() {
  if ! [[ -a "$1" ]]
    then
      ${zsb}.playMenuAndExecute cdp $__cdp_command $__cdp_dataGetter $@
    elif [[ -d "$1" ]]; then
      cd $1
    elif [[ -f "$1" ]]; then
      cd `dirname $1`
  fi
}
alias cdp='noglob __cdp'

compdef "_${zsb}.playMenuAndExecute $__cdp_dataGetter" __cdp
hisIgnore cdp

