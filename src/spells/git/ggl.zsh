ggl() {
  readonly inputBranch=${1:?'You must provide a branch to pull from.'}
  readonly currentBranch=$(${zsb}.gitBranches 'current')

  if [ "$currentBranch" != "$inputBranch" ]; then
    ${zsb}.throw "You can't use this command to pull from different branches"
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
