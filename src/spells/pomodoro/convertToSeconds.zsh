${zsb}_convertToSeconds() (
  local this="$0"
  local inputTime="$1"
  local totalSeconds

  ${this}.main () {
    ${this}.setTotalSeconds

    if [ -z "$totalSeconds" ] || ${this}.timerOutOfRange; then
      ${this}.throwBadArgumentException; return $?
    fi

    echo "$totalSeconds"
  }

  ${this}.setTotalSeconds() {
    if ${this}.isShortTimeFormat; then ${this}.setTotalSecondsFromShortTimeFormat;

    elif ${zsb}_isInteger "$inputTime"; then totalSeconds="$inputTime";

    elif ${this}.isTimeFormat; then ${this}.setTotalSecondsFromTimeFormat; fi
  }

  ${this}.isShortTimeFormat() {
    local SHORT_TIME_REGEX="^[0-9]+[hHmMsS]$"
    ${zsb}_doesMatch "$inputTime" "$SHORT_TIME_REGEX"
  }

  ${this}.isTimeFormat() {
    local TIME_REGEX="^[0-5]?[0-9]:[0-5]?[0-9](:[0-5]?[0-9])?$"
    ${zsb}_doesMatch "$inputTime" "$TIME_REGEX"
  }

  ${this}.setTotalSecondsFromShortTimeFormat() {
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

  ${this}.setTotalSecondsFromTimeFormat() {
    local reverseTime=$(echo "$inputTime" | rev)

    seconds=$(echo "$reverseTime" | cut -d":" -f1 | rev)
    minutes=$(echo "$reverseTime" | cut -d":" -f2 | rev)
    hours=$(echo "$reverseTime" | cut -d":" -f3 | rev)

    [ ! -z $hours ] && local hoursToSeconds=$(($hours * 60 * 60))
    local minutesToSeconds=$(($minutes * 60))

    totalSeconds=$(($hoursToSeconds + $minutesToSeconds + $seconds ))
  }

  ${this}.timerOutOfRange() {
    local sixtyHoursMinusOneSecond="215999"

    [ "$totalSeconds" -gt "$sixtyHoursMinusOneSecond" ] || \
      [ "$totalSeconds" -lt "1" ] && return 0
    return 1
  }

  ${this}.throwBadArgumentException() {
    echo "\
      \r${ZSB_ERROR} Bad argument.
          \rTry with (hh:)?mm:ss $(it '(min 1, max 59:59:59)')
          OR {s : s ∈ Z and 1 ≤ s ≤ 215999}
          OR ^n[hHmMsS]$ $(it '(min 1s, max 59h)')"
    return 1
  }

  ${this}.main "$@"
)

