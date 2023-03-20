multisearch() {
  case $# in
    0) ${zsb}.throw "You need at least 1 arg" ;;
    1) rg ${1} ;;
  esac

  local lastArg="${@[$#]}"
  local files=(`rg --files-with-matches ${1}`)

  (( ${#files[@]} )) || return 0

  # skip first and last arg
  for arg in ${@:2:# - 2}; do
    files=(`rg --files-with-matches ${arg} ${files[@]}`)

    (( ${#files[@]} )) || return 0
  done

  rg ${lastArg} ${files[@]}
}
