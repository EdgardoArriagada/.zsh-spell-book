# the comand is simply a [git commit --amend --no-edit && git status ]
# but it will check if the commit has already been pushed online

relabel() {
  {
    local existingFlags=(--aware --force)
    local args=("$@")
    local -A flags=( ${(z)$(${zsb}.recognizeFlags "args" "existingFlags")} )
    set -- $(${zsb}.clearFlags "args" "existingFlags")
  }

  [[ -z "$1" ]] && ${zsb}.throw "<relabel> command must contain a message."

  if ${zsb}.userWorkingOnDefaultBranch && [[ -z "${flags[--aware]}" ]]; then
    ${zsb}.throw "can't relabel in default branch, use $(hl --aware) flag to do it anyway"
  fi

  if ${zsb}.isLastCommitOnline && [[ -z "${flags[--force]}" ]]; then
    ${zsb}.throw "Can't relabel, HEAD commit has already been pushed online, use $(hl --force) flag to do it anyway."
  fi

  git commit --amend --gpg-sign --message "$*" && ${zsb}.gitStatus && echo ${zsb}.warning "Files already added to git may have been commited, use $(hl "git reset HEAD~") to undo the entire previous commit."

  return 0
}

compdef "_${zsb}.nonRepeatedListD ${ZSB_GIT_AWARE} ${ZSB_GIT_FORCE}" relabel
