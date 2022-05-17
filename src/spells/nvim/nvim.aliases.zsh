# fix [coc.nvim]
pCocNvim() {
  eval "nvm -v"
  cd ~/.config/nvim/plugged/coc.nvim && \
    eval "pyarn && yarn install && yarn build"
}
