### lazy load nvm

# for this to work, make sure you
# go to https://github.com/nvm-sh/nvm
# and follow the 'Manual Install' section

export NVM_DIR="$HOME/.nvm"

if [[ -d "$NVM_DIR" ]]; then
  local nvmScriptPath
  (( $ZSB_MACOS )) && \
    nvmScriptPath="$(brew --prefix nvm)/nvm.sh" || \
    nvmScriptPath="$NVM_DIR/nvm.sh"

  __${zsb}.prepareLazyLoad "$nvmScriptPath" \
    nvm yarn node npm vue npx tsc depcheck expo nest markdown-pdf
fi

### To run global npm packages
export PATH=~/.npm-global/bin:$PATH

