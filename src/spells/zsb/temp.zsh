# if autochdir is set in vimrc file,
# use :e temp/<new-file>.zsh to create new files

${zsb}.temp() {
  local -r callback=${1:?'You must provide a callback function.'}
  local -r tempDir=${ZSB_DIR}/src/temp

  mkdir -p ${tempDir}

  eval "${callback} ${tempDir}"
}

_${zsb}.nocompletion modifyTemp

alias cdtemp="${zsb}.temp cds"
alias vtemp="${zsb}.temp nnvim"

hisIgnore cdtemp vtemp
