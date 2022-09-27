# $1 is a pattern or a comma separated list of blobs
dondedice() {
  [[ -z "$1" ]] && return 1

  if [[ "$#" = "1" ]]
    then rg -g '!{package-lock.json}' "$1"
    else rg -g "!{package-lock.json,${1}}" "${@:2}"
  fi

}
