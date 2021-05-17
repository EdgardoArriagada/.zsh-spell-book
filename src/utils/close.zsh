close() {
  if ! xdotool -version >/dev/null 2>&1; then
    echo "${ZSB_ERROR} You need to install $(hl xdotool) first"
    return 1
  fi

  if [[ -z ${ZSB_TERMINAL_CLOSE_KEY} ]]; then
    echo " ${ZSB_WARNING} you haven't configured $(hl ZSB_TERMINAL_CLOSE_KEY) var in .env file yet"
    return 0
  fi

  xdotool key ${ZSB_TERMINAL_CLOSE_KEY}
}
