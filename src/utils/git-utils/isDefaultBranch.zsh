${zsb}.isDefaultBranch() {
  local inputBranch="$1"
  local DEFAUL_BRANCHES_REGEX="^${ZSB_GIT_DEFAULT_BRANCHES}$"
  ${zsb}.doesMatch "$inputBranch" "$DEFAUL_BRANCHES_REGEX"
}
