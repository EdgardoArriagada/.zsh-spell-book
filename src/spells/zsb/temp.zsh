# if autochdir is set in vimrc file,
# use :e temp/<new-file>.zsh to create new files

temp() {
  local tempDir=${ZSB_DIR}/src/temp

  [ ! -d ${tempDir} ] && mkdir ${tempDir}

  vim ${tempDir}
}
