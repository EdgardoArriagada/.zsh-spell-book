# send instructions to pomodoro session without attaching to it
pomodoro() {
  if ! $(tmux has-session -t pomodoro 2>/dev/null); then
    tmux new -d -s pomodoro
  fi

  local lastPomodoroLines=$(tmux capture-pane -p -t pomodoro | sed '/^$/d' | tail -1)

  if doesMatch "$lastPomodoroLines" "^⏳"; then
    echo "${ZSB_ERROR} A pomodoro timer is already running."
    return 0
  fi

  # Clean session terminal prompt
  tmux send-keys -t pomodoro.0 C-g

  # the white space at the beginning is to
  # skip it from being save to history
  local pomodoroCmd=" startPomodoro ${*}"

  tmux send-keys -t pomodoro.0 "$pomodoroCmd" ENTER && \
    echo "${ZSB_SUCCESS} Pomodoro started."
}

