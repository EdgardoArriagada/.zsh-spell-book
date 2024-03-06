# ghpr.my-repo-folder() <<< "--title 'Pr title' --body 'some body to love'"

ZSB_CREATEPR_DEFAULT_PROPS="--title WIP --body WIP"

createPr() {
  ${zsb}.assertGitRepo
  local parentBranch=${1:="$(parentBranch)"}

  if ${zsb}.userWorkingOnDefaultBranch; then
    ${zsb}.throw "Can't create a pull request on a default branch, use `hl --aware` flag to do it anyway"
  fi

  [[ -z "$parentBranch" ]] && ${zsb}.throw "Couldn't retrieve parent branch info. Pass it manually"

  ${zsb}.confirmMenu.warning "Pull request to `hl $parentBranch` from `hl $(git branch --show-current)`"

  local propsGetter="ghpr.`get_repo_name`"
  local restOfProps
  if type "$propsGetter" > /dev/null
    then restOfProps=`$propsGetter`
    else restOfProps="$ZSB_CREATEPR_DEFAULT_PROPS"
  fi

  eval "gh pr create --draft --base $parentBranch ${(z)restOfProps}"
  gh pr view --web
}

compdef "_${zsb}.singleCompC '${zsb}.gitBranches'" createPr
