tmkill() (
  local activeTmuxSessions=$(tmux ls 2>&1 | cut -d':' -s -f1)

  main() {
    if [ -z "$activeTmuxSessions" ]; then
      echo "${ZSB_INFO} There are no active tmux sessions."
      return 1
    fi

    printPrompt
    playOptionsMenu
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

  main "$@"
)
