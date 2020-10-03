dkssh() {
  if [ -z "$1" ]; then
    echo "${ZSB_ERROR} exactly one argument expected"
    return 1
  fi

  echorun docker container exec -ti "$1" bash
}

complete -C 'docker ps --format "{{.Names}}"' dkssh
