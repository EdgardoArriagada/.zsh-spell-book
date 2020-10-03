isDefaultBranch() {
  local INPUT_BRANCH="$1"
  local DEFAUL_BRANCHES_REGEX="^${ZSB_GIT_DEFAULT_BRANCHES}$"
  return $(doesMatch "$INPUT_BRANCH" "$DEFAUL_BRANCHES_REGEX")
}
