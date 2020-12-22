ggl() {
  local inputBranch="$1"
  if [ -z "$inputBranch" ]; then
    echo "${ZSB_ERROR} You must provide a branch to pull from."
    return 1
  fi

  local currentBranch=$(${zsb}.gitBranches 'current')

  if [ "$currentBranch" != "$inputBranch" ]; then
    echo "${ZSB_ERROR} You can't use this command to pull from different branches"
    return 0
  fi

  git pull origin "$inputBranch" && copythis " " -s
}
_${zsb}.ggl() {
  case $COMP_CWORD in
    1)
      COMPREPLY=( $(compgen -C "${zsb}.gitBranches 'current'") ) ;;
  esac
}

complete -F _${zsb}.ggl ggl
