copythis() (
  local GROUP_FLAGS='c'

  local inputText="$1"
  local inputFlags="$2"
  local totalArgs="$#"

  local copyCommand

  main() {
    !  areFlagsValid && return 1

    if ! isXclipInstalled; then
      throwXclipNotInstalled; return $?
    fi

    if ! areArgsValid; then
      throwInvalidUsage; return $?
    fi

    setCopyCommand
    executeCopyCommand
  }

  areFlagsValid() ${zsb}_areFlagsInGroup "$inputFlags" "$GROUP_FLAGS"

  isXclipInstalled() xclip -version >/dev/null 2>&1

  throwXclipNotInstalled () {
    echo "${ZSB_ERROR} You need to install $(hl xclip) first"
    return 1
  }

  areArgsValid() {
    [ ! -z "$inputText" ] && [ "$totalArgs" -lt "3" ] && $([ -z "$inputFlags" ] || [ "$inputFlags" = "-c" ])
  }


  throwInvalidUsage() {
    echo "${ZSB_ERROR} How to use: $(hl "copythis 'example text' -c") $(it "-c is optional")"
    return 1
  }

  setCopyCommand() copyCommand="echo -E - '$inputText' | xclip -selection clipboard"

  executeCopyCommand() {
    eval "$copyCommand"
    [[ "$inputFlags" == *'c'* ]] && return 0
    close
  }

  main "$@"
)

complete -W "-c" copythis
