copythis() {
  local GROUP_FLAGS='c'
  local flags="$2"

  ! areFlagsInGroup "$flags" "$GROUP_FLAGS" && return 1

  if ! xclip -version >/dev/null 2>&1; then
    echo "${ZSB_ERROR} You need to install $(hl xclip) first"
    return 1
  fi

  if [ -z $1 ]; then
    echo "${ZSB_ERROR} You need to provide one argument"
    return 1
  fi

  local command="echo "$1" | xclip -selection clipboard"
  if [[ "$flags" == *'c'* ]]; then
    eval "$command"
    return 0
  fi

  eval "$command" && close
}

complete -W "-c" copythis
