# usefull link https://gist.github.com/QinMing/364774610afc0e06cc223b467abe83c0#file-zshrc-L15-L37
# thanks to QinMing for the code!
${zsb}_lazy_load() {
  local o=$(timeId)
  local -a involvedScripts
  local aliasesToRelease=("${(@s: :)${1}}") # space separated list
  local scriptToLazyLoad="$2"
  local commandToRun="$3"
  shift 3
  local commandToRunArgs=$*

  {
    ${o}_main() {
      ${o}_showLoadingMessage

      ${o}_unaliasInvolvedScripts "$@"

      ${o}_sourceScriptToLazyLoad

      ${o}_runOriginalCommandWithArgs
    }

    ${o}_showLoadingMessage() {
      echo "${ZSB_INFO} Lazy loading..."
    }

    ${o}_unaliasInvolvedScripts() {
      unalias "${aliasesToRelease[@]}"
    }

    ${o}_sourceScriptToLazyLoad() {
      source "$scriptToLazyLoad"
    }

    ${o}_runOriginalCommandWithArgs() {
      eval "$commandToRun $commandToRunArgs"
    }

    ${o}_main "$@"

  } always {
    unfunction -m "${o}_*"
  }

}

_${zsb}_prepare_lazy_load() {
  local script="$1"
  shift 1

  for cmd in "$@"; do
    alias $cmd="${zsb}_lazy_load \"$*\" $script $cmd"
  done
}
