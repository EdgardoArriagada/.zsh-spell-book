Bool() {
  [[ "$1" == "0" ]] && echo true && return
  [[ -z "$1" ]] && echo false && return

  ${zsb}.isInteger "$1" && echo false && return

  [[ "$1" = "false" ]] && echo false && return

  echo true
}
