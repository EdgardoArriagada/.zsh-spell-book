${zsb}.pomodoro.runTimerFromSeconds() (
  local this="$0"
  local totalSeconds="$1"

  ${this}.main () {
    ${this}.runTimer

    ${this}.playNotifSound &
    ${this}.printEndOfTimeMsg &
    ${this}.displayNotifCard
  }

  ${this}.runTimer() {
    local -r spacesToKeepOutputClean="   "
    for i in {$totalSeconds..01}
    do
      printf "\r⏳ $(${this}.convertSecondsToTime $i)${spacesToKeepOutputClean}"
      sleep 1
    done
  }

  ${this}.convertSecondsToTime() {
    local seconds="$1"
    local hours=$(($seconds / 3600))
    local mins=$(($seconds / 60 % 60))
    local secs=$(($seconds % 60))

    [[ "$hours" -lt "10" ]] && hours="0$hours"
    [[ "$mins" -lt "10" ]] && mins="0$mins"
    [[ "$secs" -lt "10" ]] && secs="0$secs"

    if [[ "$hours" = "00" ]]; then
      if [[ "$mins" = "00" ]]; then
        echo "$secs"
        return 0
      fi
      echo "${mins}:${secs}"
      return 0
    fi
    echo "${hours}:${mins}:${secs}"
  }

  ${this}.displayNotifCard() alert "The timer for $(${this}.getCustomTimeMessage) is over"

  ${this}.getCustomTimeMessage() {
    local customTimeMessage

    if [[ "$totalSeconds" -eq "1" ]]; then
      customTimeMessage="$totalSeconds second"
    elif [[ "$totalSeconds" -lt "60" ]]; then
      customTimeMessage="$totalSeconds seconds"
    else
      customTimeMessage="$(${this}.convertSecondsToTime $totalSeconds)"
    fi

    echo "$customTimeMessage"
  }

  ${this}.playNotifSound() {
    local soundFile=${ZSB_DIR}/src/media/sounds/xylofon.wav
    [[ -f $soundFile ]] && zsb_play $soundFile
  }

  ${this}.printEndOfTimeMsg() {
    print "\r"
    ${zsb}.info "The timer for `hl "$(${this}.getCustomTimeMessage)"` was up at `hl $(date +%H:%M:%S)`"
  }

  ${this}.main "$@"
)

