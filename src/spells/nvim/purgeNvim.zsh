purgeNvim() {
  ${zsb}.confirmMenu.warning "You are about to purge nvim files"

  printAndRun "rm -rf ~/.local/share/nvim && rm -rf ~/.cache/nvim"
}
