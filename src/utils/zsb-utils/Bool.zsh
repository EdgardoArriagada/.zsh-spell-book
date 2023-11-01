Bool() {
  [[ "$1" == "0" ]] && <<< true && return
  [[ -z "$1" ]] && <<< false && return

  ${zsb}.isInteger $1 && <<< false && return

  [[ "$1" = "false" ]] && <<< false && return

  <<< true
}
