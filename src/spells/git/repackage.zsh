repackage() {
  zparseopts -D -E -F -- -aware=aware -force=force || return 1

  if ${zsb}.userWorkingOnDefaultBranch && [[ -z "$aware" ]]; then
    ${zsb}.throw "Can't repackage into default branch, use $(hl --aware) flag to do it anyway"
  fi

  if ${zsb}.isLastCommitOnline && [[ -z "$force" ]]; then
    ${zsb}.throw "Can't repackage, HEAD commit has already been pushed online, use $(hl --force) flag to do it anyway"
  fi

  git commit --amend --no-edit --gpg-sign && ${zsb}.gitStatus
  return 0
}

compdef "_${zsb}.nonRepeatedListD ${ZSB_GIT_AWARE} ${ZSB_GIT_FORCE}" repackage
