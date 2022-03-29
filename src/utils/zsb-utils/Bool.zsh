Bool() {
  [[ "$1" == "0" ]] && printf true && return
  [[ -z "$1" ]] && printf false && return

  ${zsb}.isInteger "$1" && printf false && return

  [[ "$1" = "false" ]] && printf false && return

  printf true
}
