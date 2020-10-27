${zsb}_userWorkingOnDefaultBranch() {
  local WORKING_BRANCH=$(git symbolic-ref --short HEAD)
  ${zsb}_isDefaultBranch "$WORKING_BRANCH"
}
