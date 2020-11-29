# send instructions to pomodoro tmux session without attaching to it
pomodoro() (
  local this="$0"
  local POMODORO_REGEX="^⏳"

  local totalSeconds
  local inputTime="$1"
  shift 1
  local pomodoroLabel="${*}"

  ${this}.main() {
    totalSeconds=$(${zsb}.convertToSeconds "$inputTime")

    if ! ${zsb}.didSuccess "$?"; then
      echo "$totalSeconds" # as error
      return 1
    fi

    ${this}.createPomodoroTmuxSession

    if ${this}.isPomodoroRunning; then
      ${this}.throwPomodoroAlreadyRunning; return $?
    fi

    ${this}.clearPomodoroSessionTty

    ${this}.beginPomodoro

    echo "${ZSB_SUCCESS} Pomodoro started."
    return 0
  }

  ${this}.createPomodoroTmuxSession() {
    if ! $(tmux has-session -t pomodoro 2>/dev/null); then
      tmux new -d -s pomodoro
    fi
  }

  ${this}.isPomodoroRunning() {
    local lastPomodoroLines=$(tmux capture-pane -p -t pomodoro | sed '/^$/d' | tail -1)
    ${zsb}.doesMatch "$lastPomodoroLines" "$POMODORO_REGEX"
  }

  ${this}.throwPomodoroAlreadyRunning() {
    echo "${ZSB_ERROR} A pomodoro timer is already running."
    return 1
  }

  ${this}.clearPomodoroSessionTty() tmux send-keys -t pomodoro.0 C-g

  ${this}.beginPomodoro() {
    # the white space at the beginning is to
    # skip it from being saved to zsh history
    local pomodoroCmd=" ${zsb}.startPomodoro $totalSeconds $inputTime '$pomodoroLabel'"
    tmux send-keys -t pomodoro.0 "$pomodoroCmd" ENTER
  }

  ${this}.main "$@"
)
