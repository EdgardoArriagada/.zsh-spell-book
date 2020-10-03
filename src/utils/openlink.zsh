# Open a link in your default browser
openlink() {
  local URL="$1"

  if which xdg-open >/dev/null 2>&1; then

    xdg-open "$URL" >/dev/null 2>&1

  elif which gnome-open >/dev/null 2>&1; then

    gnome-open "$URL" >/dev/null 2>&1
  else
    python -mwebbrowser "$URL" >/dev/null 2>&1
  fi
}
