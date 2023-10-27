binarieit() {
  local executable=${1?Error: You must provide an executable file.}

  local sourceBin=$(pwd)/${executable}

  sudo ln -s $(pwd)/${executable} /usr/local/bin/ &&
    ${zsb}.success "/usr/local/bin/$(hl "$1") -> $(pwdd)/$(hl ${executable})"
}
