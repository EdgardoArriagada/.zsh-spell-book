rpomodoro() {
  if ! $(tmux has-session -t pomodoro 2>/dev/null); then
    echo "${ZSB_INFO} No pomodoro session alive."
    return 0
  fi

  local lastPomodoroLines=$(tmux capture-pane -p -t pomodoro | sed '/^$/d' | tail -1)

  if ! doesMatch "$lastPomodoroLines" "^⏳"; then
    echo "${ZSB_INFO} No pomodoro timer is running."
    return 0
  fi

  echo "$lastPomodoroLines"
}
