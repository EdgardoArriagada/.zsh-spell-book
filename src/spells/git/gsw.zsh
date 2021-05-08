gsw() {
  git switch ${1}
}

# Complete with local branches excluding current
_${zsb}.gsw() {
  [ "$COMP_CWORD" -gt "1" ] && return 0

  local currentBranch=( $(${zsb}.gitBranches 'current') )
  local localBranhces=( $(${zsb}.gitBranches) )

  local newCompletionList=( )
  for branch in "${localBranhces[@]}"; do
    if [ "$branch" != "$currentBranch" ]; then
      newCompletionList+=( "$branch" )
    fi
  done

  COMPREPLY=( $(compgen -W "${newCompletionList[*]}") )
}

complete -F _${zsb}.gsw gsw
