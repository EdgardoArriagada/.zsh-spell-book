clipclean() {
  local NOTIFICATION="The clipboard has been cleaned"

  if (( $ZSB_MACOS )); then
    pbcopy <<< ' ' && say -v Karen "$NOTIFICATION"
  else
    xclip -selection clipboard <<< ' ' && \
      notify-send --urgency "low" --icon security-high "$NOTIFICATION"
  fi
}
