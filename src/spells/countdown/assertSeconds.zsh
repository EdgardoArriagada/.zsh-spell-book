${zsb}.countdown.assertSeconds() {
  local totalSeconds=$1
  local sixtyHours=216000

  (( 0 < $totalSeconds && $totalSeconds < $sixtyHours )) && return 0

  ${zsb}.throw "Bad argument.
        \rTry with (hh:)?mm:ss $(it '(min 1, max 59:59:59)')
        OR {s : s ∈ Z and 1 ≤ s ≤ 215999}
        OR ^n[hHmMsS]$ $(it '(min 1s, max 59h)')"
}
