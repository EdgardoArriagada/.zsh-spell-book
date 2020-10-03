sendmeto() (
  local inputUrl="$1"

  local inputFlags="$2"
  local GROUP_FLAGS='cs'

  main() {
    ! ${zsb}_areFlagsInGroup "$inputFlags" "$GROUP_FLAGS" && return 1

    if ! isInputUrlValid; then
      echo "${ZSB_ERROR} You must specify a valid url"
      return 1
    fi

    if inputFlagsContains "c"; then
      copyUrl && echo "${ZSB_SUCCESS} $(hl "$inputUrl") copied"
      return 0
    fi

    openlink "$inputUrl"
    didSucessOpenLink=$([ "$?" = "0" ])

    if $didSucessOpenLink && ! inputFlagsContains "s"; then
      close
    fi
  }

  isInputUrlValid() {
    [ -z "$inputUrl" ] && return 1
    local URL_REGEX="^http[s]?:\/{2}"
    return $(doesMatch "$inputUrl" "$URL_REGEX")
  }

  inputFlagsContains() return $([[ "$inputFlags" == *"$1"* ]])

  copyUrl() {
    echo "$inputUrl" | xclip -selection clipboard
  }

  main "$@"
)

complete -W "-c" sendmeto
