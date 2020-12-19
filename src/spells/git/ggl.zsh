ggl() {
  if [ -z "$1" ]; then
    echo "${ZSB_ERROR} You must provide a branch to pull from."
    return 1
  fi

  git pull origin "$1" && copythis " " -s
}

_${zsb}.ggl() {
  [ "$COMP_CWORD" -gt "1" ] && return 0

  COMPREPLY=( $(compgen -C "${zsb}.gitBranches") )
}

complete -F _${zsb}.ggl ggl
