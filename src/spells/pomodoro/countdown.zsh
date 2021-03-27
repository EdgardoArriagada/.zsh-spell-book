countdown() {
  local inputTime="$1"
  local totalSeconds="$(${zsb}.convertToSeconds "$inputTime")"

  ${zsb}.pomodoro.validateSeconds "$totalSeconds"

  ${zsb}.runTimerFromSeconds "$totalSeconds"
}

