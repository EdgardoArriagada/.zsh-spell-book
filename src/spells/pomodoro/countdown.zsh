countdown() (
  local time="$1"
  local totalSeconds

  main () {
    setTotalSeconds

    if [ -z "$totalSeconds" ] || timerOutOfRange; then
      throwBadArgumentException; return $?
    fi

    runTimer

    displayNotifCard
    playNotifSound
    printEndOfTimeMsg
  }

  setTotalSeconds() {
    if isShortTimeFormat; then setTotalSecondsFromShortTimeFormat;

    elif isInteger "$time"; then totalSeconds="$time";

    elif isTimeFormat; then setTotalSecondsFromTimeFormat; fi
  }

  isShortTimeFormat() {
    local SHORT_TIME_REGEX="^[0-9]+[hHmMsS]$"
    return $(doesMatch "$time" "$SHORT_TIME_REGEX")
  }

  isTimeFormat() {
    local TIME_REGEX="^[0-5]?[0-9]:[0-5]?[0-9](:[0-5]?[0-9])?$"
    return $(doesMatch "$time" "$TIME_REGEX")
  }

  setTotalSecondsFromShortTimeFormat() {
    local timeChar="$(echo "$time" | rev | cut -c 1)"
    local timeValue="${time%?}"

    case $timeChar in
      [hH])
        totalSeconds="$(($timeValue * 60 * 60))"
        ;;
      [mM])
        totalSeconds="$(($timeValue * 60))"
        ;;
      [sS])
        totalSeconds="$timeValue"
        ;;
      *)
        echo "${ZSB_ERROR} Unhandled error"
        ;;
    esac
  }

  setTotalSecondsFromTimeFormat() {
    local reverseTime=$(echo "$time" | rev)

    seconds=$(echo "$reverseTime" | cut -d":" -f1 | rev)
    minutes=$(echo "$reverseTime" | cut -d":" -f2 | rev)
    hours=$(echo "$reverseTime" | cut -d":" -f3 | rev)

    [ ! -z $hours ] && local hoursToSeconds=$(($hours * 60 * 60))
    local minutesToSeconds=$(($minutes * 60))

    totalSeconds=$(($hoursToSeconds + $minutesToSeconds + $seconds ))
  }

  timerOutOfRange() {
    local sixtyHoursMinusOneSecond="215999"

    [ "$totalSeconds" -gt "$sixtyHoursMinusOneSecond" ] || \
      [ "$totalSeconds" -lt "1" ] && return 0
    return 1
  }

  throwBadArgumentException() {
    echo "\
      \r${ZSB_ERROR} Bad argument.
          \rTry with (hh:)?mm:ss $(it '(min 1, max 59:59:59)')
          OR {s : s ∈ Z and 1 ≤ s ≤ 215999}
          OR ^n[hHmMsS]$ $(it '(min 1s, max 59h)')"
    return 1
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

