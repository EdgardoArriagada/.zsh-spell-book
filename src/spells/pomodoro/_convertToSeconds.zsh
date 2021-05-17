${zsb}.pomodoro.convertToSeconds() (
  local this="$0"
  local inputTime="$1"
  local totalSeconds

  ${this}.main () {
    ${this}.setTotalSeconds

    echo "$totalSeconds"
  }

  ${this}.setTotalSeconds() {
    if ${this}.isShortTimeFormat; then ${this}.setTotalSecondsFromShortTimeFormat;
    elif ${zsb}.isInteger "$inputTime"; then totalSeconds="$inputTime";
    elif ${this}.isTimeFormat; then ${this}.setTotalSecondsFromTimeFormat;
    else totalSeconds=0;
    fi
  }

  ${this}.isShortTimeFormat() {
    local SHORT_TIME_REGEX="^[0-9]+[hHmMsS]$"
    ${zsb}.doesMatch "$inputTime" "$SHORT_TIME_REGEX"
  }

  ${this}.isTimeFormat() {
    local TIME_REGEX="^[0-5]?[0-9]:[0-5]?[0-9](:[0-5]?[0-9])?$"
    ${zsb}.doesMatch "$inputTime" "$TIME_REGEX"
  }

  ${this}.setTotalSecondsFromShortTimeFormat() {
    local timeChar="$(echo "$inputTime" | rev | cut -c 1)"
    local timeValue="${inputTime%?}"

    case $timeChar in
      [hH])
        totalSeconds="$(($timeValue * 60 * 60))" ;;
      [mM])
        totalSeconds="$(($timeValue * 60))" ;;
      [sS])
        totalSeconds="$timeValue" ;;
      *)
        echo "${ZSB_ERROR} Unhandled error" ;;
    esac
  }

  ${this}.setTotalSecondsFromTimeFormat() {
    local -r reverseTime=$(echo "$inputTime" | rev)

    local -r seconds=$(echo "$reverseTime" | cut -d":" -f1 | rev)
    local -r minutes=$(echo "$reverseTime" | cut -d":" -f2 | rev)
    local -r hours=$(echo "$reverseTime" | cut -d":" -f3 | rev)

    local hoursToSeconds=0
    [[ ! -z "$hours" ]] && hoursToSeconds=$(($hours * 60 * 60))
    local -r minutesToSeconds=$(($minutes * 60))

    totalSeconds=$(($hoursToSeconds + $minutesToSeconds + $seconds ))
  }

  ${this}.main "$@"
)

