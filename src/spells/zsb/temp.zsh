# if autochdir is set in vimrc file,
# use :e temp/<new-file>.zsh to create new files

${zsb}.temp() {
  local callback=${1:?'You must provide a callback function.'}
  local tempDir=${ZSB_DIR}/src/temp

  mkdir -p ${tempDir}

  ${callback} ${tempDir}
}

_${zsb}.nocompletion modifyTemp

alias cdtemp="${zsb}.temp cds"
alias vtemp="${zsb}.temp nvim"

hisIgnore cdtemp vtemp
