Bool() {
  local input="$1"

  if ${zsb}.isFloat "$input"; then
    [[ "$1" = "0" ]] && echo true && return
    echo false && return
  fi

  [[ -z "$input" ]] && echo false && return

  echo true
}
