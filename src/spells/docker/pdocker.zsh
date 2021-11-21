pdocker() {
  docker --version >/dev/null 2>&1 && ${zsb}.throw "Docker is already installed."

  ${zsb}.confirmMenu.withPrompt && curl -sSL https://get.docker.com/ | sh

  # To be able to run docker without sudo
  sudo groupadd docker
  sudo usermod -aG docker $USER

  ${zsb}.warning "Log out and log back in and kill your tmux server to be able to use docker without root privileges."
}
