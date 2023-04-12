${zsb}.doesBranchesMatch() {
  local inputBranch=${1}
  local branchesRegex=${2}
  ${zsb}.doesMatch ${inputBranch} "^${branchesRegex}$"
}

${zsb}.isDefaultBranch() ${zsb}.doesBranchesMatch ${1} ${ZSB_GIT_DEFAULT_BRANCHES}
${zsb}.isCriticalBranch() ${zsb}.doesBranchesMatch ${1} ${ZSB_GIT_CRITICAL_BRANCHES}

${zsb}.userWorkingOnDefaultBranch() ${zsb}.isDefaultBranch `git branch --show-current`
${zsb}.userWorkingOnCriticalBranch() ${zsb}.isCriticalBranch `git branch --show-current`

