tmkill() (
  local this="$0"
  local activeTmuxSessions=$(tmls)
  local inputSessions=( "$@" )

  ${this}.main() {
    if ! ${this}.existsActiveTmuxSessions; then
      ${this}.throwNoTmuxSessions; return $?
    fi

    if ${this}.inputSessionsReceived; then
      ${this}.killInputSessions; return $?
    fi

    ${this}.playMenu
  }

  ${this}.inputSessionsReceived() [[ "${#inputSessions[@]}" -gt "0" ]]

  ${this}.existsActiveTmuxSessions() [[ -n "$activeTmuxSessions" ]]

  ${this}.killInputSessions() {
    for session in "${inputSessions[@]}"; do
      tmux kill-session -t "$session" >/dev/null 2>&1
    done
  }

  ${this}.playMenu() {
    ${this}.printPrompt
    ${this}.killTmuxServerOnConfirm
  }

  ${this}.printPrompt() {
    ${zsb}.warning "The following tmux sessions will be deleted:"
    echo " "
    echo "$(hl "$activeTmuxSessions")"
    echo " "
    ${zsb}.prompt "Do you really want to delete these sessions? [Y/n]."
  }

  ${this}.killTmuxServerOnConfirm() {
    ${zsb}.confirmMenu && ${this}.killTmuxServer &&
      ${zsb}.success "All tmux sessions have been deleted."
  }

  ${this}.killTmuxServer() tmux kill-server

  ${this}.throwNoTmuxSessions() {
    ${zsb}.info "There are no active tmux sessions."
    return 0
  }

  ${this}.main "$@"
)

compdef "_${zsb}.nonRepeatedListC tmls" tmkill

