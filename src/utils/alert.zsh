alert() {
  zparseopts -D -E -F -- i:=_icon || return 1
  local icon="${_icon[2]}"
  : ${icon:='timer'}

  (( $ZSB_MACOS )) && \
    osascript -e "display notification \"$*\"" || \
    notify-send --icon "$icon" --urgency critical "$*"
}

alias a="alert"
alias ad="alert 'DONE'"

