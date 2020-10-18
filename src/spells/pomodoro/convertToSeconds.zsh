${zsb}_convertToSeconds() (
  local inputTime="$1"
  local totalSeconds

  main () {
    setTotalSeconds

    if [ -z "$totalSeconds" ] || timerOutOfRange; then
      throwBadArgumentException; return $?
    fi

    echo "$totalSeconds"
  }

  setTotalSeconds() {
    if isShortTimeFormat; then setTotalSecondsFromShortTimeFormat;

    elif ${zsb}_isInteger "$inputTime"; then totalSeconds="$inputTime";

    elif isTimeFormat; then setTotalSecondsFromTimeFormat; fi
  }

  isShortTimeFormat() {
    local SHORT_TIME_REGEX="^[0-9]+[hHmMsS]$"
    return $(${zsb}_doesMatch "$inputTime" "$SHORT_TIME_REGEX")
  }

  isTimeFormat() {
    local TIME_REGEX="^[0-5]?[0-9]:[0-5]?[0-9](:[0-5]?[0-9])?$"
    return $(${zsb}_doesMatch "$inputTime" "$TIME_REGEX")
  }

  setTotalSecondsFromShortTimeFormat() {
    local timeChar="$(echo "$inputTime" | rev | cut -c 1)"
    local timeValue="${inputTime%?}"

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
    local reverseTime=$(echo "$inputTime" | rev)

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

  main "$@"
)

