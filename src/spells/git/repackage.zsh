repackage() (
  zparseopts -D -E -F -- -aware=aware -force=force -no-verify=noVerify || return 1
  local this="$0"

  ${this}.validateNonDefaultBranch() {
    ${zsb}.userWorkingOnDefaultBranch &&
      ${zsb}.throw "can't repackage in default branch, use `hl --aware` flag to do it anyway"
  }

  ${this}.validateCommitOffline() {
    ${zsb}.isLastCommitOnline &&
      ${zsb}.throw "Can't repackage, HEAD commit has already been pushed online, use `hl --force` flag to do it anyway."
  }

  ${this}.warnUnintendedAddedFiles() {
    ${zsb}.warning "Files already added to git may have been commited, use `hl 'git reset HEAD~'` to undo the entire previous commit."
  }

  [[ -z "$aware" ]] && ${this}.validateNonDefaultBranch

  [[ -z "$force" ]] && ${this}.validateCommitOffline

  if [[ -n "$1" ]]; then
    amend --message "$*" ${noVerify} &&
      ${zsb}.gitStatus &&
      ${this}.warnUnintendedAddedFiles
  else
    amend --no-edit ${noVerify} &&
      ${zsb}.gitStatus
  fi
)

hisIgnore repackage

_${zsb}.repackage() {
  local noVerify
  local aware

  ${zsb}.isHusky && noVerify='--no-verify:Skip husky verifications'
  ${zsb}.userWorkingOnDefaultBranch && aware=${ZSB_GIT_AWARE}

  _${zsb}.nonRepeatedListD ${aware} ${ZSB_GIT_FORCE} ${noVerify}
}

compdef _${zsb}.repackage repackage
