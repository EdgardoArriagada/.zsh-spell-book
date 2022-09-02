alert() {
  zparseopts -D -E -F -- i:=_icon || return 1
  local icon="${_icon[2]}"
  : ${icon:='timer'}

  if (( $ZSB_MACOS ))
    then osascript -e "tell app \"System Events\" to display dialog \"${*}\"" >/dev/null 2>&1
    else notify-send --icon "$icon" --urgency critical "$*"
  fi
}

alias a="alert"
alias ad="alert 'DONE'"

