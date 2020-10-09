tmkill() (
  local activeTmuxSessions=$(tml)
  local inputSession="$1"

  main() {
    if ! existsActiveTmuxSessions; then
      throwNoTmuxSessions; return $?
    fi

    if tooManyArgs "$@"; then
      throwTooManyArgsException; return $?
    fi

    if existsInputSession; then
      killInputSession; return $?
    fi

    printPrompt
    playOptionsMenu
  }

  tooManyArgs() return $([ "$#" -gt "1" ])

  existsActiveTmuxSessions() return $([ ! -z "$activeTmuxSessions" ])

  existsInputSession() return $([ ! -z "$inputSession" ])

  killInputSession() {
      return $(tmux kill-session -t "$inputSession")
  }

  printPrompt() {
    echo "${ZSB_WARNING} The following tmux sessions will be deleted:"
    echo " "
    echo "$(hl "$activeTmuxSessions")"
    echo " "
    echo "${ZSB_PROMPT} Do you really want to delete these sessions? [Y/n]."
  }

  playOptionsMenu() {
    ${zsb}_yesNoMenu killTmuxServer &&
      echo "${ZSB_SUCCESS} All tmux sessions have been deleted."
  }

  killTmuxServer() tmux kill-server

  throwNoTmuxSessions() {
    echo "${ZSB_INFO} There are no active tmux sessions."
    return 0
  }

  throwTooManyArgsException() {
    echo "${ZSB_ERROR} Only one argument expected."
    return 1
  }

  main "$@"
)

complete -C "tml" tmkill
