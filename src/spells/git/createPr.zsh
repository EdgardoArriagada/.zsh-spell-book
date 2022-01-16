# ghpr.my-repo-folder() printf "--title 'Pr title' --body 'some body to love'"

createPr_defaultProps="--title WIP --body WIP"

createPr() {
  ${zsb}.validateGitRepo
  zparseopts -D -E -F -- -aware=aware || return 1

  local parentBranch=`gitParentBranch`
  [[ -z "$parentBranch" ]] && ${zsb}.cancel "There is not a parent branch."

  if ${zsb}.userWorkingOnDefaultBranch && [[ -z "$aware" ]]; then
    ${zsb}.throw "Can't create a pull request on a default branch, use `hl --aware` flag to do it anyway"
  fi

  ${zsb}.warning "Pull request to `hl ${parentBranch}` from `hl $(git branch --show-current)`"
  ${zsb}.confirmMenu.withPrompt

  local propsGetter="ghpr.`getRepoName`"
  local restOfProps
  if type "$propsGetter" > /dev/null; then
    restOfProps=`${propsGetter}`
  else
    restOfProps="$createPr_defaultProps"
  fi

  gh pr create --draft --base "$parentBranch" ${(z)restOfProps}
  gh pr view --web
}

compdef "_${zsb}.nonRepeatedList --aware" createPr
