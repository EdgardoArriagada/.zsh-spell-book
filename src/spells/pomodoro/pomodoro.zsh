# send instructions to pomodoro tmux session without attaching to it
pomodoro() (
  local POMODORO_REGEX="^⏳"

  local totalSeconds
  local inputTime="$1"
  shift 1
  local pomodoroLabel="${*}"

  main() {
    totalSeconds=$(${zsb}_convertToSeconds "$inputTime")

    if ! ${zsb}_didSuccess "$?"; then
      echo "$totalSeconds" # as error
      return 1
    fi

    createPomodoroTmuxSession

    if isPomodoroRunning; then
      throwPomodoroAlreadyRunning; return $?
    fi

    clearPomodoroSessionTty

    beginPomodoro

    echo "${ZSB_SUCCESS} Pomodoro started."
    return 0
  }

  createPomodoroTmuxSession() {
    if ! $(tmux has-session -t pomodoro 2>/dev/null); then
      tmux new -d -s pomodoro
    fi
  }

  isPomodoroRunning() {
    local lastPomodoroLines=$(tmux capture-pane -p -t pomodoro | sed '/^$/d' | tail -1)
    return $(${zsb}_doesMatch "$lastPomodoroLines" "$POMODORO_REGEX")
  }

  throwPomodoroAlreadyRunning() {
    echo "${ZSB_ERROR} A pomodoro timer is already running."
    return 1
  }

  clearPomodoroSessionTty() tmux send-keys -t pomodoro.0 C-g

  beginPomodoro() {
    # the white space at the beginning is to
    # skip it from being saved to zsh history
    local pomodoroCmd=" ${zsb}_startPomodoro $totalSeconds $inputTime '$pomodoroLabel'"
    tmux send-keys -t pomodoro.0 "$pomodoroCmd" ENTER
  }

  main "$@"
)
