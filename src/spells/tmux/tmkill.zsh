tmkill() (
  local this="$0"
  local activeTmuxSessions=$(tmls)
  local inputSessions=( "$@" )

  ${this}.inputSessionsReceived() [[ "${#inputSessions[@]}" -gt "0" ]]

  ${this}.existsActiveTmuxSessions() [[ -n "$activeTmuxSessions" ]]

  ${this}.killInputSessions() {
    for session in "${inputSessions[@]}"; do
      tmux kill-session -t "$session" >/dev/null 2>&1
    done
  }

  ${this}.playMenu() {
    ${zsb}.confirmMenu.withItems \
      "The following tmux sessions will be deleted:" \
      "$activeTmuxSessions"

    tmux kill-server &&
      ${zsb}.success "All tmux sessions have been deleted."
  }

  ${this}.validateTmuxSessions() {
    if ! ${this}.existsActiveTmuxSessions; then
      ${zsb}.cancel "There are no active tmux sessions."
    fi
  } 

  ${this}.validateTmuxSessions

  if ${this}.inputSessionsReceived; then
    ${this}.killInputSessions; return $?
  fi

  ${this}.playMenu
)

compdef "_${zsb}.nonRepeatedListC tmls" tmkill

