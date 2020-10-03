userWorkingOnDefaultBranch() {
  local WORKING_BRANCH=$(git symbolic-ref --short HEAD)
  return $(isDefaultBranch "$WORKING_BRANCH")
}
