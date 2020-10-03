startPomodoro() (
  local time="$1"
  shift 1
  local label="${*}"
  local beginString
  local thisMonthFolder

  main() {
    if [ -z "$time" ]; then
      echo "${ZSB_ERROR} 1 arg expected."; return 1
    fi

    setBeginString

    setThisMonthFolder

    appendSpaceToLog

    appendPomodoroInfoToLog "started"

    countdown "$time"

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
    logToZsb "${beginString} $time $verb at $(generateDate)" "$thisMonthFolder"
  }

  generateDate() echo $(date +%H:%M:%S)

  setThisMonthFolder() {
    local year="$(date +%Y)"
    local month="$(date +%b)"
    thisMonthFolder="pomodoro/${year}/${month}"
  }

  main "$@"
)
