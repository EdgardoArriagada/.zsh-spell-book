messageStatus() {
  local msg="\e[38;5;244m${@}\033[0m"

  "$@"

  (( $? )) && \
    printf "\r\e[38;5;1m✘\033[0m ${msg}\n" && \
    return 0

  printf "\r\e[38;5;2m✔\033[0m ${msg}\n"
}
