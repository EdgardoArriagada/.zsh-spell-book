sendmeto() (
  local this="$0"
  local inputUrl="$1"
  local inputFlags="$2"
  local GROUP_FLAGS='cs'

  ${this}.main() {
    ! ${this}.areFlagsValid && return $?

    if ! ${this}.isInputUrlValid; then
      ${this}.throwInvalidUrl; return $?
    fi

    if ${this}.inputFlagsContains "c"; then
      ${this}.copyInputUrlToClipboard;
    else
      ${this}.openInputUrlInBrowser
    fi

    ${this}.inputFlagsContains "s" && return $?

    close
  }

  ${this}.areFlagsValid() ${zsb}_areFlagsInGroup "$inputFlags" "$GROUP_FLAGS"

  ${this}.isInputUrlValid() {
    local URL_REGEX="^http[s]?:\/{2}"
    ${zsb}_doesMatch "$inputUrl" "$URL_REGEX"
  }

  ${this}.throwInvalidUrl() {
    echo "${ZSB_ERROR} You must specify a valid url"
    return 1
  }

  ${this}.inputFlagsContains() [[ "$inputFlags" == *"$1"* ]]

  ${this}.copyInputUrlToClipboard() {
    ${this}.copyUrl && echo "${ZSB_SUCCESS} $(hl "$inputUrl") copied"
    return 0
  }

  ${this}.openInputUrlInBrowser() openlink "$inputUrl"

  ${this}.copyUrl() {
    echo "$inputUrl" | xclip -selection clipboard
  }

  ${this}.main "$@"
)

complete -W "-c" sendmeto
