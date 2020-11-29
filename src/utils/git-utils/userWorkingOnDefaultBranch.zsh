${zsb}.userWorkingOnDefaultBranch() {
  local WORKING_BRANCH=$(git symbolic-ref --short HEAD)
  ${zsb}.isDefaultBranch "$WORKING_BRANCH"
}
