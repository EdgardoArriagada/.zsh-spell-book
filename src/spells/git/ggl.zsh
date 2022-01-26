ggl() {
  local currentBranch=`git_current_branch`

  ${zsb}.info "Pulling from `hl ${currentBranch}`"

  git pull origin "$currentBranch"
}

_${zsb}.nocompletion ggl

