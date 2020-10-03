echorun() {
  echo "${ZSB_INFO} Running $(hl "${*}")"
  eval "$@"
}
