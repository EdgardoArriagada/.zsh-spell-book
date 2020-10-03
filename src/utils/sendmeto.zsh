sendmeto() (
  local argumentURL="$1"

  local flags="$2"
  local GROUP_FLAGS='cs'

  local highlightedURL="$(hl "${argumentURL}")"

  main() {
    ! areFlagsInGroup "$flags" "$GROUP_FLAGS" && return 1

    if ! isUrl "$argumentURL" || [ -z $argumentURL ]; then
      echo "${ZSB_ERROR} You must specify a valid url"
      return 1
    fi

    if [[ "$flags" == *'c'* ]]; then
      copyUrl && echo "${ZSB_SUCCESS} ${highlightedURL} copied"
      return 0
    fi

    openlink "$argumentURL"
    didSucessOpenLink=$([ "$?" = "0" ])


    if $didSucessOpenLink && [[ "$flags" != *'s'* ]]; then
      close
    fi

  }

  isUrl() {
    local URL_REGEX="^http[s]?:\/{2}"
    return $(doesMatch "$1" "$URL_REGEX")
  }

  copyUrl() {
    echo "$argumentURL" | xclip -selection clipboard
  }

  main "$@"
)

complete -W "-c" sendmeto
