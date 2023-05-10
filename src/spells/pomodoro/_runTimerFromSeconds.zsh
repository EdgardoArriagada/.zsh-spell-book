${zsb}.pomodoro.runTimerFromSeconds() (
  local this=$0
  local totalSeconds=$1

  ${this}.runTimer() {
    local -r spacesToKeepOutputClean='   '
    for i in {${totalSeconds}..01}; do
      printf "\r⏳ `${this}.convertSecondsToTime ${i}`${spacesToKeepOutputClean}"
      sleep 1
    done
  }

  # Refactored to perform the less amount of instructions
  ${this}.convertSecondsToTime() {
    # $1 input total seconds
    local hours=$(($1 / 3600))
    local mins=$(($1 / 60 % 60))
    local secs=$(($1 % 60))

    (( $secs < 10 )) && secs=0${secs}

    if (( $hours == 0 )); then
      if (( $mins == 0 ))
        then printf ${secs}
        else (( $mins < 10 )) && mins=0${mins}; printf ${mins}:${secs}
      fi
      return 0
    fi

    (( ${hours} < 10 )) && hours=0${hours}
    (( ${mins} < 10 )) && mins=0${mins}

    printf ${hours}:${mins}:${secs}
  }

  ${this}.displayNotifCard() alert "The timer for `${this}.getCustomTimeMessage` is over"

  ${this}.getCustomTimeMessage() {
    local customTimeMessage

    if (( ${totalSeconds} == 1 ))
      then customTimeMessage="${totalSeconds} second"

    elif (( ${totalSeconds} < 60 ))
      then customTimeMessage="${totalSeconds} seconds"

    else
      customTimeMessage="`${this}.convertSecondsToTime ${totalSeconds}`"
    fi

    printf ${customTimeMessage}
  }

  ${this}.playNotifSound() {
    local soundFile=${ZSB_DIR}/src/media/sounds/xylofon.wav
    [[ -f ${soundFile} ]] && zsb_play ${soundFile}
  }

  ${this}.printEndOfTimeMsg() {
    print "\r"
    ${zsb}.info "The timer for `hl $(${this}.getCustomTimeMessage)` was up at `hl $(date +%H:%M:%S)`"
  }

   { # main
    ${this}.runTimer

    ${this}.playNotifSound &
    ${this}.printEndOfTimeMsg &
    ${this}.displayNotifCard
   }
)

