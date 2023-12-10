iconize.fmt() printf "\e[38;5;244m${*}\033[0m"

iconize.output() {
  $@

  if (( $? ))
    then printf "\r\e[38;5;1m✘\033[0m `iconize.fmt $@`\n"
    else printf "\r\e[38;5;2m✔\033[0m `iconize.fmt $@`\n"
  fi
}

iconize.skip() {
  printf "\r\e[38;5;3mskip\033[0m `iconize.fmt $@`\n"
}

iconize.notFound() {
  printf "\r\e[38;5;6mnot found\033[0m `iconize.fmt $@`\n"
}
