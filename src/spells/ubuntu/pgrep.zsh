pgrep() {
  [ "$#" -lt "2" ] && echo "${ZSB_ERROR} 2 or more args expected" && return 1

  local toGrep="${@: -1}" # last argument
  set -- "${@:0:$#}" # remove last argument
  shift 1 # remove "pgrep" from list of args

  local resultCommand="${*} | grep '$toGrep'"
  echo "${ZSB_INFO} Running $(hl "$resultCommand")"
  eval "$resultCommand"
}

alias pg='pgrep'
