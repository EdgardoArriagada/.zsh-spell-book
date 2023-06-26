# usefull link https://gist.github.com/QinMing/364774610afc0e06cc223b467abe83c0#file-zshrc-L15-L37
# thanks to QinMing for the code!
${zsb}.lazyLoad() {
  local aliasesToRelease=("${(@s: :)${1}}") # space separated list
  local scriptToLazyLoad=$2
  local commandToRun=$3
  shift 3
  local commandToRunArgs=$*

  ${zsb}.info "Lazy loading..."

  unalias ${aliasesToRelease[@]}
  source ${scriptToLazyLoad}
  eval "${commandToRun} ${commandToRunArgs}"
}

${zsb}.prepareLazyLoad() {
  local script=$1
  shift 1

  for cmd in ${@}; do
    alias ${cmd}="${zsb}.lazyLoad \"${*}\" ${script} ${cmd}"
  done
}
