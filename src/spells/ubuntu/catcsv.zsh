catcsv() {
  column -t -s, -n "$@" | less -F -S -X -K
}
