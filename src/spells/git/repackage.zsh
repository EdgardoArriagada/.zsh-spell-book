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

  [[ -z "$noVerify" ]] && ${zsb}.activateNvmIfHusky

  if [[ -n "$1" ]]; then
    git commit --amend --gpg-sign --message "$*" ${noVerify} &&
      ${zsb}.gitStatus &&
      ${this}.warnUnintendedAddedFiles
  else
    git commit --amend --gpg-sign --no-edit ${noVerify} &&
      ${zsb}.gitStatus
  fi
)

hisIgnore repackage

_${zsb}.repackage() {
  local noVerify
  ${zsb}.isHusky && noVerify='--no-verify:Skip husky verifications'

  _${zsb}.nonRepeatedListD "$ZSB_GIT_AWARE" "$ZSB_GIT_FORCE" "$noVerify"
}

compdef _${zsb}.repackage repackage
