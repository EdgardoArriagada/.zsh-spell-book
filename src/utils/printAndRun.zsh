printAndRun() {
  local inputCommand="$1"
  echo -n "${ZSB_INFO} "
  plainPrint "Running $(hl "$inputCommand")"
  eval "$inputCommand"
}
