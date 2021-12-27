binarieit() {
  local executable=${1:?'You must provide an executable file'}

  sudo ln -s $(pwd)/${1} /usr/local/bin/ &&
    ${zsb}.success "/usr/local/bin/$(hl "$1")"
}
