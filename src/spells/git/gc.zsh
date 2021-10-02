gc() {
  zparseopts -D -E -F -- -aware=aware || return 1

  if ${zsb}.userWorkingOnDefaultBranch && [[ -z "$aware" ]]; then
    ${zsb}.throw "Can't commit into default branch, use $(hl "--aware") flag to do it anyway"
  fi

  local commitSuccess=false
  if [[ -z "$1" ]]; then
    git commit --gpg-sign && ${zsb}.gitStatus && commitSuccess=true
  else
    git commit --gpg-sign -m "$*" && ${zsb}.gitStatus && commitSuccess=true
  fi

  if ${zsb}.userWorkingOnDefaultBranch && "$commitSuccess"; then
    ${zsb}.warning "Commit made into default branch, use $(hl "git reset HEAD~") to undo the entire previous commit"
  fi

  return 0
}

compdef "_${zsb}.nonRepeatedList --aware" gc

