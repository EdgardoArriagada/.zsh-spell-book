dcossh() {
  if [ -z "$1" ]; then
    echo "${ZSB_ERROR} exactly one argument expected"
    return 1
  fi

  echorun docker-compose exec "$1" bash
}

complete -C 'docker-compose ps --services' dcossh
