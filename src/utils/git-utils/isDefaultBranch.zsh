${zsb}.isDefaultBranch() {
  local inputBranch=${1}
  local DEFAUL_BRANCHES_REGEX="^${ZSB_GIT_DEFAULT_BRANCHES}$"
  ${zsb}.doesMatch ${inputBranch} ${DEFAUL_BRANCHES_REGEX}
}

${zsb}.isCriticalBranch() {
  local inputBranch=${1}
  local CRITICAL_BRANCHES_REGEX="^${ZSB_GIT_CRITICAL_BRANCHES}$"
  ${zsb}.doesMatch ${inputBranch} ${CRITICAL_BRANCHES_REGEX}
}
