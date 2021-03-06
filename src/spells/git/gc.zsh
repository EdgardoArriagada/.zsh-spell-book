# commit changes when not on DEFAULT_BRANCH branch and asks for status after
# passing a message as string is not necessary, but a best practice (tested on ubuntu 18.04 with zsh)

gc() {
  {
    declare -A args
    args[--aware]=false
    ${zsb}.switchTrueMatching "${args[@]}" "$@"
    set -- $(${zsb}.clearFlags "${args[@]}" "$@")
  }

  if ${zsb}.userWorkingOnDefaultBranch && ! "${args[--aware]}"; then
    echo "${ZSB_ERROR} Can't commit into default branch, use $(hl "--aware") flag to do it anyway"
    return 1
  fi

  local commitSuccess=false
  if [[ -z "$1" ]]; then
    git commit --gpg-sign && ${zsb}.gitStatus && commitSuccess=true
  else
    git commit --gpg-sign -m "$*" && ${zsb}.gitStatus && commitSuccess=true
  fi

  if ${zsb}.userWorkingOnDefaultBranch && "$commitSuccess"; then
    echo "${ZSB_WARNING} Commit made into default branch, use $(hl "git reset HEAD~") to undo the entire previous commit"
  fi

  return 0
}

compdef "_${zsb}.nonRepeatedList --aware" gc

