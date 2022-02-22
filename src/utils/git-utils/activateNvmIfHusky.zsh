${zsb}.isHusky() {
  [[ -f .huskyrc ]] || [[ -d .husky ]]
}

${zsb}.activateNvmIfHusky() {
  ${zsb}.isHusky && eval 'nvm --version'
}
