tmkill() (
  local this="$0"
  local activeTmuxSessions=$(tml)
  local inputSession="$1"

  ${this}.main() {
    if ! ${this}.existsActiveTmuxSessions; then
      ${this}.throwNoTmuxSessions; return $?
    fi

    if ${this}.tooManyArgs "$@"; then
      ${this}.throwTooManyArgsException; return $?
    fi

    if ${this}.existsInputSession; then
      ${this}.killInputSession; return $?
    fi

    ${this}.printPrompt
    ${this}.playOptionsMenu
  }

  ${this}.tooManyArgs() [ "$#" -gt "1" ]

  ${this}.existsActiveTmuxSessions() [ ! -z "$activeTmuxSessions" ]

  ${this}.existsInputSession() [ ! -z "$inputSession" ]

  ${this}.killInputSession() {
     tmux kill-session -t "$inputSession"
  }

  ${this}.printPrompt() {
    echo "${ZSB_WARNING} The following tmux sessions will be deleted:"
    echo " "
    echo "$(hl "$activeTmuxSessions")"
    echo " "
    echo "${ZSB_PROMPT} Do you really want to delete these sessions? [Y/n]."
  }

  ${this}.playOptionsMenu() {
    ${zsb}_yesNoMenu ${this}.killTmuxServer &&
      echo "${ZSB_SUCCESS} All tmux sessions have been deleted."
  }

  ${this}.killTmuxServer() tmux kill-server

  ${this}.throwNoTmuxSessions() {
    echo "${ZSB_INFO} There are no active tmux sessions."
    return 0
  }

  ${this}.throwTooManyArgsException() {
    echo "${ZSB_ERROR} Only one argument expected."
    return 1
  }

  ${this}.main "$@"
)

complete -C "tml" tmkill
