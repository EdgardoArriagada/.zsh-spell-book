dkssh() {
  if [ -z "$1" ]; then
    echo "${ZSB_ERROR} exactly one argument expected"
    return 1
  fi

  printAndRun "docker container exec -ti $1 bash"
}

_${zsb}.dkssh() {
  [ "$COMP_CWORD" -gt "1" ] && return 0

  COMPREPLY=( $(compgen -C "docker ps --format "{{.Names}}"") )
}

complete -F _${zsb}.dkssh dkssh
