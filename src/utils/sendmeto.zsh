sendmeto() (
  local inputUrl="$1"
  local inputFlags="$2"
  local GROUP_FLAGS='cs'

  main() {
    ! areFlagsValid && return $?

    if ! isInputUrlValid; then
      throwInvalidUrl; return $?
    fi

    if inputFlagsContains "c"; then
      copyInputUrlToClipboard;
    else
      openInputUrlInBrowser
    fi

    inputFlagsContains "s" && return $?

    close
  }

  areFlagsValid() ${zsb}_areFlagsInGroup "$inputFlags" "$GROUP_FLAGS"

  isInputUrlValid() {
    [ -z "$inputUrl" ] && return 1
    local URL_REGEX="^http[s]?:\/{2}"
    return $(${zsb}_doesMatch "$inputUrl" "$URL_REGEX")
  }

  throwInvalidUrl() {
    echo "${ZSB_ERROR} You must specify a valid url"
    return 1
  }

  inputFlagsContains() return $([[ "$inputFlags" == *"$1"* ]])

  copyInputUrlToClipboard() {
    copyUrl && echo "${ZSB_SUCCESS} $(hl "$inputUrl") copied"
    return 0
  }

  openInputUrlInBrowser() openlink "$inputUrl"

  copyUrl() {
    echo "$inputUrl" | xclip -selection clipboard
  }

  main "$@"
)

complete -W "-c" sendmeto
