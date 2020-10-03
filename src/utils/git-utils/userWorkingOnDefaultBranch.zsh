userWorkingOnDefaultBranch() {
  local WORKING_BRANCH=$(git symbolic-ref --short HEAD)
  return $(${zsb}_isDefaultBranch "$WORKING_BRANCH")
}
