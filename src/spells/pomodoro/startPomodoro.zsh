${zsb}_startPomodoro() (
  local totalSeconds="$1"
  local inputTime="$2"
  local label="$3"
  local beginString
  local thisMonthFolder

  main() {
    setBeginString

    setThisMonthFolder

    appendSpaceToLog

    appendPomodoroInfoToLog "started"

    ${zsb}_runTimerFromSeconds "$totalSeconds"

    appendPomodoroInfoToLog "ended"
  }

  setBeginString() {
    if [ -z "$label" ]; then
      beginString="Session for"
      return 0
    fi
    local uppercaseLabel="${(C)label}"
    beginString="${uppercaseLabel} for"
  }

  appendSpaceToLog() logToZsb " " "$thisMonthFolder"

  appendPomodoroInfoToLog() {
    local verb="$1" # stared | ended
    logToZsb "${beginString} $inputTime $verb at $(generateDate)" "$thisMonthFolder"
  }

  generateDate() echo $(date +%H:%M:%S)

  setThisMonthFolder() {
    local year="$(date +%Y)"
    local month="$(date +%b)"
    thisMonthFolder="pomodoro/${year}/${month}"
  }

  main "$@"
)
