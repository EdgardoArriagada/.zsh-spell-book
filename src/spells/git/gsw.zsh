gsw() {
  git switch ${1}
}

_${zsb}.gsw() {
  [ "$CURRENT" -gt "2" ] && return 0

  local -a subcmds
  local localBranhces=( $(${zsb}.gitBranches) )
  local currentBranch=( $(${zsb}.gitBranches 'current') )

  subcmds=( ${localBranhces:|currentBranch} )
  _describe 'command' subcmds
}

compdef _${zsb}.gsw gsw

