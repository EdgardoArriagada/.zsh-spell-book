multisearch() {
  case $# in
    0) ${zsb}.throw "You need at least 1 arg" ;;
    1) rg ${1} ;;
  esac

  local lastArg="${@[$#]}"
  local files=(`rg --files-with-matches ${1}`)

  for arg in ${@:2}; do
    [[ "$arg" = "$lastArg" ]] && break
    (( ${#files[@]} )) || break

    files=(`rg --files-with-matches ${arg} $(print -r -- -g ${^files})`)
  done

  (( ${#files[@]} )) || return 0

  rg ${lastArg} `print -r -- -g ${^files}`
}


