cdr() {
  local choosenDir=$(dirs -p | fzf)
  [[ -z "$choosenDir" ]] && return 0
  eval "cd ${choosenDir}"
}
