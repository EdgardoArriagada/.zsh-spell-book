dondedice() {
  [[ -z "$1" ]] && return 1

  if [[ "$#" == "1" ]]; then
    rg -g '!{package-lock.json}' "$1"
    return $? 
  fi

  rg -g "!{package-lock.json,${1}}" "${@:2}"
}
