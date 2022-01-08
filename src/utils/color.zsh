color() {
  if [[ "$1" = "-h" ]]; then
    print "\nUsage: \"color `hl x` my message\", where `hl x` is a number in the following list:\n"
    printColours
    return 0
  fi

  local inputColor="$1"
  shift 1

  printf "\x1b[38;5;${inputColor}m${*}\x1b[0m\n"
}
