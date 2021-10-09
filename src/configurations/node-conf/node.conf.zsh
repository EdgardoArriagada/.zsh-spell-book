### lazy load nvm

# for this to work, make sure you
# go to https://github.com/nvm-sh/nvm
# and follow the 'Manual Install' section

export NVM_DIR="$HOME/.nvm"

if [[ -s "$NVM_DIR/nvm.sh" ]] && [[ -s "$NVM_DIR/bash_completion" ]]; then
  __${zsb}.prepareLazyLoad "$NVM_DIR/nvm.sh" \
    nvm yarn node npm vue npx tsc depcheck expo nest markdown-pdf
fi

### To run global npm packages
export PATH=~/.npm-global/bin:$PATH

