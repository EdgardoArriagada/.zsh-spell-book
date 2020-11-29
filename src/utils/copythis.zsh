copythis() (
  local this="$0"
  local GROUP_FLAGS='c'

  local inputText="$1"
  local inputFlags="$2"
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

    ${this}.setCopyCommand
    ${this}.executeCopyCommand
  }

  ${this}.areFlagsValid() ${zsb}.areFlagsInGroup "$inputFlags" "$GROUP_FLAGS"

  ${this}.isXclipInstalled() xclip -version >/dev/null 2>&1

  ${this}.throwXclipNotInstalled () {
    echo "${ZSB_ERROR} You need to install $(hl xclip) first"
    return 1
  }

  ${this}.areArgsValid() {
    [ ! -z "$inputText" ] && [ "$totalArgs" -lt "3" ] && $([ -z "$inputFlags" ] || [ "$inputFlags" = "-c" ])
  }


  ${this}.throwInvalidUsage() {
    echo "${ZSB_ERROR} How to use: $(hl "copythis 'example text' -c") $(it "-c is optional")"
    return 1
  }

  ${this}.setCopyCommand() copyCommand="echo -E - '$inputText' | xclip -selection clipboard"

  ${this}.executeCopyCommand() {
    eval "$copyCommand"
    [[ "$inputFlags" == *'c'* ]] && return 0
    close
  }

  ${this}.main "$@"
)

complete -W "-c" copythis
