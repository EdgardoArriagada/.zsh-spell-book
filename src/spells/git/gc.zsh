gc() {
  zparseopts -D -E -F -- -aware=aware || return 1

  if ${zsb}.userWorkingOnDefaultBranch && [[ -z "$aware" ]]; then
    ${zsb}.throw "Can't commit into default branch, use $(hl "--aware") flag to do it anyway"
  fi

  if [[ -z "$1" ]]; then
    git commit --gpg-sign
  else
    git commit --gpg-sign -m "$*"
  fi

  local didCommitSuccess="$(Bool $?)"

  "$didCommitSuccess" || return 1

  ${zsb}.gitStatus

  if ${zsb}.userWorkingOnDefaultBranch; then
    ${zsb}.warning "Commit made into default branch, use $(hl "git reset HEAD~") to undo the entire previous commit"
  fi
}

compdef "_${zsb}.nonRepeatedList --aware" gc

