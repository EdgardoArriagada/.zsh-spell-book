gc() {
  zparseopts -D -E -F -- -aware=aware -no-verify=noVerify || return 1

  ${zsb}.userWorkingOnDefaultBranch
  local isUserInDefaultBranch="$(Bool $?)"

  if "$isUserInDefaultBranch" && [[ -z "$aware" ]]; then
    ${zsb}.throw "Can't commit into default branch, use `hl --aware` flag to do it anyway"
  fi

  if [[ -z "$1" ]]
    then git commit --gpg-sign ${noVerify} || return 1
    else git commit --gpg-sign -m "$*" ${noVerify} || return 1
  fi

  ${zsb}.gitStatus

  if "$isUserInDefaultBranch"; then
    ${zsb}.warning "Commit made into default branch, use `hl git reset HEAD~` to undo the entire previous commit"
  fi
}

_${zsb}.gc() {
  local noVerify
  local aware

  ${zsb}.isHusky && noVerify='--no-verify:Skip husky verifications'
  ${zsb}.userWorkingOnDefaultBranch && aware=${ZSB_GIT_AWARE}

  _${zsb}.nonRepeatedListD ${aware} ${noVerify}
}

compdef _${zsb}.gc gc

