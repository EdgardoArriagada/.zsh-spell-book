# usefull link https://gist.github.com/QinMing/364774610afc0e06cc223b467abe83c0#file-zshrc-L15-L37
# thanks to QinMing for the code!
${zsb}_lazy_load() {
  local this="$0_$(${zsb}_timeId)"
  local -a involvedScripts
  local aliasesToRelease=("${(@s: :)${1}}") # space separated list
  local scriptToLazyLoad="$2"
  local commandToRun="$3"
  shift 3
  local commandToRunArgs=$*

  {
    ${this}.main() {
      ${this}.showLoadingMessage

      ${this}.unaliasInvolvedScripts "$@"

      ${this}.sourceScriptToLazyLoad

      ${this}.runOriginalCommandWithArgs
    }

    ${this}.showLoadingMessage() {
      echo "${ZSB_INFO} Lazy loading..."
    }

    ${this}.unaliasInvolvedScripts() {
      unalias "${aliasesToRelease[@]}"
    }

    ${this}.sourceScriptToLazyLoad() {
      source "$scriptToLazyLoad"
    }

    ${this}.runOriginalCommandWithArgs() {
      eval "$commandToRun $commandToRunArgs"
    }

    ${this}.main "$@"
  } always {
    unfunction -m "${this}.*"
  }
}

__${zsb}_prepare_lazy_load() {
  local script="$1"
  shift 1

  for cmd in "$@"; do
    alias $cmd="${zsb}_lazy_load \"$*\" $script $cmd"
  done
}
