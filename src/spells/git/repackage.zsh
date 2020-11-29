# the comand is simply a [git commit --amend --no-edit --gpg-sign && git status ]
# but it will check if the commit has already been pushed online

repackage() {
  {
    declare -A args
    args[--aware]=false
    args[--force]=false
    ${zsb}.switchTrueMatching "${args[@]}" "$@"
  }

  if ${zsb}.userWorkingOnDefaultBranch && ! "${args[--aware]}"; then
    echo "${ZSB_ERROR} Can't repackage into default branch, use ${ZSB_SHL}--aware${ZSB_EHL} flag to do it anyway"
    return 1
  fi

  if ${zsb}.isLastCommitOnline && ! "${args[--force]}"; then
    echo "${ZSB_ERROR} Can't repackage, HEAD commit has already been pushed online, use ${ZSB_SHL}--force${ZSB_EHL} flag to do it anyway"
    return 1
  fi

  git commit --amend --no-edit --gpg-sign && ${zsb}.gitStatus
  return 0
}

_${zsb}.repackage() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList=( "--aware" "--force" )
  local newCompletion=( $(${zsb}.removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )

  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}
complete -F _${zsb}.repackage repackage
