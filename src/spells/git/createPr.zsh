# ghpr.my-repo-folder() printf "--title 'Pr title' --body 'some body to love'"

createPr() {
  zparseopts -D -E -F -- -aware=aware || return 1

  local currDir=${PWD##*/}
  local propsGetter="ghpr.${currDir}"
  ! type "$propsGetter" > /dev/null && ${zsb}.throw "there is not a props getter for `hl ${currDir}`"

  local parentBranch=`gitParentBranch`
  [[ -z "$parentBranch" ]] && ${zsb}.cancel "There is not a parentBranch."

  if ${zsb}.userWorkingOnDefaultBranch && [[ -z "$aware" ]]; then
    ${zsb}.throw "Can't create a pull request on a default branch, use `hl --aware` flag to do it anyway"
  fi

  ${zsb}.warning "Pull request to `hl ${parentBranch}` from `hl $(git branch --show-current)`"
  ${zsb}.confirmMenu.withPrompt

  local restOfProps=`${propsGetter}`
  gh pr create --draft --base "$parentBranch" ${(z)restOfProps}
  gh pr view --web
}

compdef "_${zsb}.nonRepeatedList --aware" createPr
