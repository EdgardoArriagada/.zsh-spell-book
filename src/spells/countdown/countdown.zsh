countdown() {
  local this=${zsb}.${0}
  local inputTime=$1
  local totalSeconds=`${this}.convertToSeconds ${inputTime}`

  ${this}.validateSeconds ${totalSeconds}

  ${this}.runTimerFromSeconds ${totalSeconds}
}

