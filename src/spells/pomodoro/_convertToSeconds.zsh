${zsb}.pomodoro.convertToSeconds() (
  local this=$0
  local inputTime=$1
  local totalSeconds

  ${this}.setTotalSeconds() {
    if ${this}.isShortTimeFormat
      then ${this}.setTotalSecondsFromShortTimeFormat

    elif ${zsb}.isInteger ${inputTime}
      then totalSeconds=${inputTime}

    elif ${this}.isTimeFormat
      then ${this}.setTotalSecondsFromTimeFormat

    else totalSeconds=0
    fi
  }

  ${this}.isShortTimeFormat() {
    local SHORT_TIME_REGEX='^[0-9]+[hHmMsS]$'
    ${zsb}.doesMatch ${inputTime} ${SHORT_TIME_REGEX}
  }

  ${this}.isTimeFormat() {
    local TIME_REGEX='^[0-5]?[0-9]:[0-5]?[0-9](:[0-5]?[0-9])?$'
    ${zsb}.doesMatch ${inputTime} ${TIME_REGEX}
  }

  ${this}.setTotalSecondsFromShortTimeFormat() {
    local timeChar=`rev <<< ${inputTime} | cut -c 1`
    local timeValue=${inputTime%?}

    case $timeChar in
      [hH])
        totalSeconds=$((${timeValue} * 60 * 60)) ;;
      [mM])
        totalSeconds=$((${timeValue} * 60)) ;;
      [sS])
        totalSeconds=$timeValue ;;
      *)
        ${zsb}.throw 'Unhandled error' ;;
    esac
  }

  ${this}.setTotalSecondsFromTimeFormat() {
    local reverseTime=`rev <<< ${inputTime}`

    local seconds=`cut -d':' -f1 <<< ${reverseTime} | rev`
    local minutes=`cut -d':' -f2 <<< ${reverseTime} | rev`
    local hours=`cut -d':' -f3 <<< ${reverseTime} | rev`

    local hoursToSeconds=0
    [[ -n "$hours" ]] && hoursToSeconds=$((${hours} * 60 * 60))
    local minutesToSeconds=$((${minutes} * 60))

    totalSeconds=$((${hoursToSeconds} + ${minutesToSeconds} + ${seconds} ))
  }

  { # main
    ${this}.setTotalSeconds

    <<< ${totalSeconds}
  }
)

