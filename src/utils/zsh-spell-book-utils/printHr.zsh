${zsb}.printHr() {
  printf %"$COLUMNS"s | tr " " "_"
  [ "$1" = "--no-space" ] && return 0
  echo " "
}
