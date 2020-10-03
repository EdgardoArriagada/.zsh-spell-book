resource() {
  clear

  if [ ! "$(echo $ZSH_VERSION)" ]; then
    source ~/.bashrc
    return 0
  fi

  source ~/.zshrc
}
