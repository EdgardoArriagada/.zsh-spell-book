# the comand is simply a [git commit --amend --no-edit && git status ]
# but it will check if the commit has already been pushed online

relabel() {
  zparseopts -D -E -F -- -aware=aware -force=force

  [[ -z "$1" ]] && ${zsb}.throw "<relabel> command must contain a message."

  if ${zsb}.userWorkingOnDefaultBranch && [[ -z "$aware" ]]; then
    ${zsb}.throw "can't relabel in default branch, use $(hl --aware) flag to do it anyway"
  fi

  if ${zsb}.isLastCommitOnline && ! [[ -z "$force" ]]; then
    ${zsb}.throw "Can't relabel, HEAD commit has already been pushed online, use $(hl --force) flag to do it anyway."
  fi

  git commit --amend --gpg-sign --message "$*" && ${zsb}.gitStatus &&
    ${zsb}.warning "Files already added to git may have been commited, use $(hl "git reset HEAD~") to undo the entire previous commit."

  return 0
}

compdef "_${zsb}.nonRepeatedListD ${ZSB_GIT_AWARE} ${ZSB_GIT_FORCE}" relabel
