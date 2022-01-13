# ${zsb}.ghpr_my-repo-folder() printf "--title -WIP --base main"

createPr() {
  local currDir=${PWD##*/}
  local propsGetter="${zsb}.ghpr_${currDir}"
  ! type "$propsGetter" > /dev/null && ${zsb}.throw "there is not a props getter for `hl ${currDir}`"

  ${zsb}.confirmMenu.withPrompt

  local prProps=`${propsGetter}`
  gh pr create --draft ${(z)prProps}
  gh pr view --web
}

alias gprcreate='createPr'
alias GPRCREATE='toggleCapsLock && gprcreate'
