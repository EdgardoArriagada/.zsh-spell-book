countdown() {
  local inputTime=$1
  local totalSeconds=`${zsb}.pomodoro.convertToSeconds ${inputTime}`

  ${zsb}.pomodoro.validateSeconds ${totalSeconds}

  ${zsb}.pomodoro.runTimerFromSeconds ${totalSeconds}
}

