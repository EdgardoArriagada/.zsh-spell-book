# ghpr.my-repo-folder() printf "--title 'Pr title' --body 'some body to love'"

createPr() {
  local currDir=${PWD##*/}
  local propsGetter="ghpr.${currDir}"
  ! type "$propsGetter" > /dev/null && ${zsb}.throw "there is not a props getter for `hl ${currDir}`"

  local parentBranch=`gitParentBranch`
  [[ -z "$parentBranch" ]] && ${zsb}.cancel "There is not a parentBranch."

  ${zsb}.warning "Pull request to `hl ${parentBranch}` from `hl $(git branch --show-current)`"
  ${zsb}.confirmMenu.withPrompt

  local restOfProps=`${propsGetter}`
  gh pr create --draft --base "$parentBranch" ${(z)restOfProps}
  gh pr view --web
}
