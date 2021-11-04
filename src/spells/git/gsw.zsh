gsw() {
  git switch ${1}
}

_${zsb}.gsw() {
  [[ "$CURRENT" -gt "2" ]] && return 0

  local localBranhces=( $(${zsb}.gitBranches) )
  local currentBranch=( $(${zsb}.gitBranches 'current') )

  local newCompletion=( ${localBranhces:|currentBranch} )
  _${zsb}.verticalComp "newCompletion"
}

compdef _${zsb}.gsw gsw

