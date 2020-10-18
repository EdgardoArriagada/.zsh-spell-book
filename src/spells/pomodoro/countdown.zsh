countdown() {
  local inputTime="$1"
  totalSeconds=$(${zsb}_convertToSeconds "$inputTime")

  if ! ${zsb}_didSuccess "$?"; then
    echo "$totalSeconds" # as error
    return 1
  fi

  ${zsb}_runTimerFromSeconds "$totalSeconds"
}
