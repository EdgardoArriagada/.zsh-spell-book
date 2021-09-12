${zsb}.notebook.validateEnv() {
  ${zsb}.validate "ZSB_NOTEBOOK_DIR"
  ${zsb}.validate "ZSB_NOTEBOOK_CHAPTER"
  ${zsb}.validate "ZSB_NOTEBOOK_PAGE"
}

${zsb}.notebook() {
  local -r callback=${1:?'You must provide a callback'}

  ${zsb}.notebook.validateEnv

  local -r currentPage="${ZSB_NOTEBOOK_DIR}/${ZSB_NOTEBOOK_CHAPTER}/${ZSB_NOTEBOOK_PAGE}"

  # Create parent folder
  mkdir -p $(dirname $currentPage)

  eval "${callback} ${currentPage}"
}

alias cdpage="${zsb}.notebook.validateEnv && cds ${ZSB_NOTEBOOK_DIR}/${ZSB_NOTEBOOK_CHAPTER}"
alias page="${zsb}.notebook nnvim"
alias cpage="${zsb}.notebook batcat"
alias ccppage="${zsb}.notebook ccp"

