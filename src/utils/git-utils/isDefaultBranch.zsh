${zsb}_isDefaultBranch() {
  local INPUT_BRANCH="$1"
  local DEFAUL_BRANCHES_REGEX="^${ZSB_GIT_DEFAULT_BRANCHES}$"
  return $(${zsb}_doesMatch "$INPUT_BRANCH" "$DEFAUL_BRANCHES_REGEX")
}
