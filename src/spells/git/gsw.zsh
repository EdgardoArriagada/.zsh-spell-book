gsw() {
  git switch ${1}
}

# Complete with local branches excluding current
_${zsb}.gsw() {
  [ "$COMP_CWORD" -gt "1" ] && return 0

  local currentBranch="$(${zsb}.gitBranches 'current')"
  local localBranhces=( $(${zsb}.gitBranches) )

  for branch in "${localBranhces[@]}"; do
    [ "$branch" = "$currentBranch" ] && continue
    COMPREPLY+=( "$branch" )
  done
}

complete -F _${zsb}.gsw gsw
