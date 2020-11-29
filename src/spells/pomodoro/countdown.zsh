countdown() {
  local inputTime="$1"
  totalSeconds=$(${zsb}.convertToSeconds "$inputTime")

  if ! ${zsb}.didSuccess "$?"; then
    echo "$totalSeconds" # as error
    return 1
  fi

  ${zsb}.runTimerFromSeconds "$totalSeconds"
}
