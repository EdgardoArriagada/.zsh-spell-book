${zsb}.doesBranchesMatch() {
  local inputBranch=$1
  local branchesRegex=$2

  [[ $inputBranch =~ "^${branchesRegex}$" ]]
}

${zsb}.isDefaultBranch() ${zsb}.doesBranchesMatch $1 $ZSB_GIT_DEFAULT_BRANCHES
${zsb}.isMainBranch() ${zsb}.doesBranchesMatch $1 $ZSB_GIT_MAIN_BRANCHES

${zsb}.isUserOnDefaultBranch() ${zsb}.isDefaultBranch `git branch --show-current`
${zsb}.isUserOnMainBranch() ${zsb}.isMainBranch `git branch --show-current`

