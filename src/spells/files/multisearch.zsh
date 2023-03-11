multisearch() {
  case $# in
    0) ${zsb}.throw "You need at least 1 arg" ;;
    1) rg ${1} ;;
  esac

  local lastArg="${@[$#]}"
  local files=(`rg --files-with-matches ${1}`)

  # skip first and last arg
  for arg in ${@:2:# - 2}; do
    (( ${#files[@]} )) || break

    files=(`rg --files-with-matches ${arg} ${files[@]}`)
  done

  (( ${#files[@]} )) || return 0

  rg ${lastArg} ${files[@]}
}


