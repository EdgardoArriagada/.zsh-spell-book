${zsb}.pomodoro.startPomodoro() (
  local this=$0
  local totalSeconds=$1
  local inputTime=$2
  local label=$3
  local beginString
  local thisMonthFolder

  ${this}.setBeginString() {
    if [[ -z "$label" ]]
      then beginString='Session for'
      else beginString="${label} for"
    fi
  }

  ${this}.setThisMonthFolder() {
    local year=`date +%Y`
    local month=`date +%b`
    thisMonthFolder=pomodoro/${year}/${month}
  }

  ${this}.appendSpaceToLog() logToZsb ' ' ${thisMonthFolder}

  ${this}.appendPomodoroInfoToLog() {
    local verb=$1 # stared | ended
    logToZsb "${beginString} ${inputTime} ${verb} at `${this}.generateDate`" ${thisMonthFolder}
  }

  ${this}.generateDate() date +%H:%M:%S

  { # main
    ${this}.setBeginString

    ${this}.setThisMonthFolder

    ${this}.appendSpaceToLog

    ${this}.appendPomodoroInfoToLog 'started'

    ${zsb}.pomodoro.runTimerFromSeconds "$totalSeconds"

    ${this}.appendPomodoroInfoToLog 'ended'
  }
)
