${zsb}.searchFileAndOpen() {
  local choosenDir=$(fd -t f | fzf)
  [[ -z "$choosenDir" ]] && return 0
  eval "${1} ${choosenDir}"
}

alias cg="${zsb}.searchFileAndOpen c"
alias vg="${zsb}.searchFileAndOpen v"
