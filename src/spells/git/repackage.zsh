# the comand is simply a [git commit --amend --no-edit --gpg-sign && git status ]
# but it will check if the commit has already been pushed online

repackage() {
  {
    local existingFlags=(--aware --force)
    local args=("$@")
    local -A flags=( ${(z)$(${zsb}.recognizeFlags "args" "existingFlags")} )
  }

  if ${zsb}.userWorkingOnDefaultBranch && ! "${flags[--aware]}"; then
    ${zsb}.throw "Can't repackage into default branch, use $(hl --aware) flag to do it anyway"
  fi

  if ${zsb}.isLastCommitOnline && ! "${flags[--force]}"; then
    ${zsb}.throw "Can't repackage, HEAD commit has already been pushed online, use $(hl --force) flag to do it anyway"
  fi

  git commit --amend --no-edit --gpg-sign && ${zsb}.gitStatus
  return 0
}

compdef "_${zsb}.nonRepeatedListD ${ZSB_GIT_AWARE} ${ZSB_GIT_FORCE}" repackage
