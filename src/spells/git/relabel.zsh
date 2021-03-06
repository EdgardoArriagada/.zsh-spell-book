# the comand is simply a [git commit --amend --no-edit && git status ]
# but it will check if the commit has already been pushed online

relabel() {
  {
    declare -A args
    args[--aware]=false
    args[--force]=false
    ${zsb}.switchTrueMatching "${args[@]}" "$@"
    set -- $(${zsb}.clearFlags "${args[@]}" "$@")
  }

  if [[ -z "$1" ]]; then
    echo "${ZSB_ERROR} <relabel> command must contain a message"
    return 1
  fi

  if ${zsb}.userWorkingOnDefaultBranch && ! "${args[--aware]}"; then
    echo "${ZSB_ERROR} can't relabel in default branch, use ${ZSB_SHL}--aware${ZSB_EHL} flag to do it anyway"
    return 1
  fi

  if ${zsb}.isLastCommitOnline && ! "${args[--force]}"; then
    echo "${ZSB_ERROR} Can't relabel, HEAD commit has already been pushed online, use ${ZSB_SHL}--force${ZSB_EHL} flag to do it anyway"
    return 1
  fi

  git commit --amend --gpg-sign --message "$*" && ${zsb}.gitStatus && echo "${ZSB_WARNING} files already added to git may have been commited, use ${ZSB_SHL}git reset HEAD~${ZSB_EHL} to undo the entire previous commit"

  return 0
}

compdef "_${zsb}.nonRepeatedListD ${ZSB_GIT_AWARE} ${ZSB_GIT_FORCE}" relabel
