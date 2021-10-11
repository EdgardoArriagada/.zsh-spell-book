alert() {
  zparseopts -D -E -F -- i:=_icon || return 1
  local icon="${_icon[2]}"
  : ${icon:='timer'}

  notify-send --icon "$icon" --urgency critical "$*"
}

alias a="alert"
alias ad="alert 'DONE'"

