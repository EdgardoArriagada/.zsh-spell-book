gsw() {
  case $1 in
    [0=]) git switch - ;;
    *) git switch ${1} ;;
  esac
}

hisIgnore 'gsw -' 'gsw 0' 'gsw ='

_${zsb}.gsw() {
  [[ "$CURRENT" -gt "2" ]] && return 0

  local localBranhces=( $(${zsb}.gitBranches) )
  local currentBranch=( $(${zsb}.gitBranches 'current') )

  local newCompletion=( ${localBranhces:|currentBranch} )
  _describe 'command' newCompletion
}

compdef _${zsb}.gsw gsw

