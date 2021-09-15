clipclean() {
  print " " |
    xclip -selection clipboard &&
    notify-send --urgency "low" --icon security-high "The clipboard has been cleaned"
}
