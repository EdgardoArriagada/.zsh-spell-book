### lazy load nvm

# for this to work, make sure you
# go to https://github.com/nvm-sh/nvm
# and follow the 'Manual Install' section

export NVM_DIR="$HOME/.nvm"

${zsb}.lazyNvm() {
  (( $ZSB_MACOS )) && \
    source "$(brew --prefix nvm)/nvm.sh" || \
    source "$NVM_DIR/nvm.sh"

  ${zsb}.removeNodeVersionCacheDecorator nvm
}

if [[ -d "$NVM_DIR" ]]; then
  local nvmScriptPath

  __${zsb}.prepareLazyLoad "${zsb}.lazyNvm" \
    nvm yarn node npm vue npx tsc depcheck expo nest markdown-pdf
fi

### To run global npm packages
export PATH=~/.npm-global/bin:$PATH

