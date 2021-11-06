${zsb}.page.validateEnv() {
  ${zsb}.validate "ZSB_NOTEBOOK_DIR"
  ${zsb}.validate "ZSB_NOTEBOOK_CHAPTER"
}

${zsb}.page.getCurrentPage() {


  local year=$(date +%Y)
  local month=$(date +%b)
  local page=$(date +%b-%d-%Y.md)

  local currentPage="${ZSB_NOTEBOOK_DIR}/${ZSB_NOTEBOOK_CHAPTER}/${year}/${month}/${page}"
  mkdir -p $(dirname $currentPage)

  print $currentPage
}

${zsb}.page() {${zsb}.page.validateEnv && eval "$@ $(${zsb}.page.getCurrentPage)"}

cdpage() {${zsb}.page.validateEnv && cds $(dirname $(${zsb}.page.getCurrentPage))}

alias cdnotebook="${zsb}.validate 'ZSB_NOTEBOOK_DIR' && cds ${ZSB_NOTEBOOK_DIR}"
alias page="(cdpage && ${zsb}.page nnvim)"
alias cpage="${zsb}.page batcat"
alias ccppage="${zsb}.page ccp"

