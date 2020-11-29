dcossh() {
  if [ -z "$1" ]; then
    echo "${ZSB_ERROR} exactly one argument expected"
    return 1
  fi

  printAndRun "docker-compose exec $1 bash"
}

_${zsb}.dcossh() {
  [ "$COMP_CWORD" -gt "1" ] && return 0

  COMPREPLY=( $(compgen -C "docker-compose ps --services") )
}

complete -F _${zsb}.dcossh dcossh
