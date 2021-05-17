copythis() (
  local this="$0"
  local GROUP_FLAGS='ps'

  local inputText=$(${zsb}.nonPiped "$1")
  local inputFlags=$(${zsb}.nonPiped "$2" "$1")

  local totalArgs="$#"

  local copyCommand


  ${this}.main() {
    !  ${this}.areFlagsValid && return 1

    if ! ${this}.isXclipInstalled; then
      ${this}.throwXclipNotInstalled; return $?
    fi

    if ! ${this}.areArgsValid; then
      ${this}.throwInvalidUsage; return $?
    fi

    ${this}.executeCopyCommand
    ${this}.executeSideEffects
  }

  ${this}.areFlagsValid() ${zsb}.areFlagsInGroup "$inputFlags" "$GROUP_FLAGS"

  ${this}.isXclipInstalled() xclip -version >/dev/null 2>&1

  ${this}.throwXclipNotInstalled () {
    echo "${ZSB_ERROR} You need to install $(hl xclip) first"
    return 1
  }

  ${this}.areArgsValid() {
    [[ ! -z "$inputText" ]] && [[ "$totalArgs" -lt "3" ]] && $([[ -z "$inputFlags" ]] || [[ "$inputFlags" == '-'* ]])
  }


  ${this}.throwInvalidUsage() {
    echo "${ZSB_ERROR} How to use: $(hl "copythis 'example text' -<flag>") $(it "-<flag> is optional")"
    return 1
  }

  ${this}.executeCopyCommand() {
    echo -En - ${inputText} | xclip -selection clipboard
  }

  ${this}.executeSideEffects() {
    [[ "$inputFlags" == *'p'* ]] && echo "${ZSB_SUCCESS} $(hl ${inputText}) copied!"
    [[ "$inputFlags" != *'s'* ]] && close
  }

  ${this}.main "$@"
)

complete -W "-ps" copythis
