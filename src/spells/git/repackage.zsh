repackage() {
  zparseopts -D -E -F -- -aware=aware -force=force || return 1

  if ${zsb}.userWorkingOnDefaultBranch && [[ -z "$aware" ]]; then
    ${zsb}.throw "can't repackage in default branch, use $(hl --aware) flag to do it anyway"
  fi

  if ${zsb}.isLastCommitOnline && ! [[ -z "$force" ]]; then
    ${zsb}.throw "Can't repackage, HEAD commit has already been pushed online, use $(hl --force) flag to do it anyway."
  fi

  if [[ -n "$1" ]]; then
    git commit --amend --gpg-sign --message "$*" && ${zsb}.gitStatus &&
      ${zsb}.warning "Files already added to git may have been commited, use $(hl "git reset HEAD~") to undo the entire previous commit."
  else
    git commit --amend --gpg-sign --message && ${zsb}.gitStatus
  fi
}

hisIgnore repackage

compdef "_${zsb}.nonRepeatedListD '${ZSB_GIT_AWARE}' '${ZSB_GIT_FORCE}'" repackage
