purgeNvim() {
  ${zsb}.confirmMenu.warning "You are about to purge nvim files"

  spinner rm -rf ~/.local/share/nvim
  spinner rm -rf ~/.local/state/nvim
  spinner rm -rf ~/.cache/nvim
}
