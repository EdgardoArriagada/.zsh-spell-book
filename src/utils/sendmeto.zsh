sendmeto() (
  local this="$0"
  local inputUrl
  local stayInScreen
  local copyToClipboard

  ${this}.main() {
    zparseopts -D -E -F -- s=stayInScreen c=copyToClipboard || return 1
    inputUrl="${@[1]}"

    if ! ${zsb}.isUrl "$inputUrl"; then
      ${this}.throwInvalidUrl
    fi

    if [[ -n "$copyToClipboard" ]]; then
      ${this}.copyInputUrlToClipboard;
    else
      ${this}.openInputUrlInBrowser
    fi

    ${this}.executeSideEffects
  }

  ${this}.throwInvalidUrl() ${zsb}.throw "You must specify a valid url"

  ${this}.inputFlagsContains() [[ "$inputFlags" == *"$1"* ]]

  ${this}.copyInputUrlToClipboard() {
    ${this}.copyUrl && echo "${ZSB_SUCCESS} $(hl "$inputUrl") copied"
    return 0
  }

  ${this}.openInputUrlInBrowser() openlink "$inputUrl"

  ${this}.copyUrl() {
    echo "$inputUrl" | xclip -selection clipboard
  }

  ${this}.executeSideEffects() {
    [[ -z "$stayInScreen" ]] && close
  }

  ${this}.main "$@"
)

complete -W "-c" sendmeto
