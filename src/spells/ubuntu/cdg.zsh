cdg() {
  local choosenDir=$(fd -t d | fzf)
  [[ -z "$choosenDir" ]] && return 0
  eval "cd ${choosenDir}"
}
