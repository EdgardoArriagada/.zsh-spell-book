${zsb}_runTimerFromSeconds() (
  local totalSeconds="$1"

  main () {
    runTimer

    displayNotifCard
    playNotifSound
    printEndOfTimeMsg
  }

  runTimer() {
    local spacesToKeepOutputClean="   "
    for i in {$totalSeconds..01}
    do
      echo -ne "\r⏳ $(convertSecondsToTime $i)${spacesToKeepOutputClean}"
      sleep 1
    done
  }

  convertSecondsToTime() {
    local seconds="$1"
    local hours=$(($seconds / 3600))
    local mins=$(($seconds / 60 % 60))
    local secs=$(($seconds % 60))

    [ "$hours" -lt "10" ] && hours="0$hours"
    [ "$mins" -lt "10" ] && mins="0$mins"
    [ "$secs" -lt "10" ] && secs="0$secs"

    if [ "$hours" = "00" ]; then
      if [ "$mins" = "00" ]; then
        echo "$secs"
        return 0
      fi
      echo "$mins:$secs"
      return 0
    fi
    echo "$hours:$mins:$secs"
  }

  displayNotifCard() alert "The timer for $(getCustomTimeMessage) is over"

  getCustomTimeMessage() {
    local customTimeMessage

    if [ "$totalSeconds" -eq "1" ]; then
      customTimeMessage="$totalSeconds second"
    elif [ "$totalSeconds" -lt "60" ]; then
      customTimeMessage="$totalSeconds seconds"
    else
      customTimeMessage="$(convertSecondsToTime $totalSeconds)"
    fi

    echo "$customTimeMessage"
  }

  playNotifSound() {
    local soundFile=${ZSB_DIR}/src/media/sounds/xylofon.wav
    [ -f $soundFile ] && aplay $soundFile>/dev/null 2>&1
  }

  printEndOfTimeMsg() echo "\r${ZSB_INFO} The timer for $(hl "$(getCustomTimeMessage)") was up at $(hl $(date +%H:%M:%S))"

  main "$@"
)

