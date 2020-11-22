ggl() {
  if [ -z "$1" ]; then
    echo "${ZSB_ERROR} You must provide a branch to pull from."
    return 1
  fi

  git pull origin "$1" && copythis " " -c
}

_${zsb}_ggl() {
  [ "$COMP_CWORD" -gt "1" ] && return 0

  COMPREPLY=( $(compgen -C "${zsb}_gitBranches") )
}

complete -F _${zsb}_ggl ggl
