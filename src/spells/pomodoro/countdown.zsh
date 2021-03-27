countdown() (
  local this="$0"
  local inputTime="$1"
  local totalSeconds

  ${this}.main () {
    totalSeconds="$(${zsb}.convertToSeconds "$inputTime")"

    if ${this}.isTimerOutOfRange; then
      ${this}.throwBadArgumentException
    fi

    ${zsb}.runTimerFromSeconds "$totalSeconds"
  }


  ${this}.isTimerOutOfRange() {
    local sixtyHoursMinusOneSecond="215999"

    [ "$totalSeconds" -gt "$sixtyHoursMinusOneSecond" ] || \
    [ "$totalSeconds" -lt "1" ] && return 0
    return 1
  }

  ${this}.throwBadArgumentException() {
    ${zsb}.throw "\
      \r${ZSB_ERROR} Bad argument.
          \rTry with (hh:)?mm:ss $(it '(min 1, max 59:59:59)')
          OR {s : s ∈ Z and 1 ≤ s ≤ 215999}
          OR ^n[hHmMsS]$ $(it '(min 1s, max 59h)')"
  }

  ${this}.main "$@"
)




