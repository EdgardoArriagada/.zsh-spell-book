${zsb}.page.validateEnv() {
  ${zsb}.validate "ZSB_NOTEBOOK_DIR"
  ${zsb}.validate "ZSB_NOTEBOOK_CHAPTER"
}

${zsb}.page.getTodaysPage() {
  local year=$(date +%Y)
  local month=$(date +%b)
  local page=$(date +%b-%d-%Y.md)

  local currentPage="${ZSB_NOTEBOOK_DIR}/${ZSB_NOTEBOOK_CHAPTER}/${year}/${month}/${page}"
  mkdir -p $(dirname ${currentPage})

  print ${currentPage}
}

alias cdnotebook="${zsb}.validate 'ZSB_NOTEBOOK_DIR' && cds ${ZSB_NOTEBOOK_DIR}"

${zsb}.page() {
  local cmd="$1"
  local currentPage="$2"
  ${zsb}.page.validateEnv

  if [[ -z "$currentPage" ]]; then
    local todaysPage=`${zsb}.page.getTodaysPage`

    (builtin cd $(dirname ${todaysPage}) && \
      ${cmd} ${todaysPage})
  else
    (builtin cd ${ZSB_NOTEBOOK_DIR}/${ZSB_NOTEBOOK_CHAPTER} && \
      ${cmd} "$currentPage")
  fi
}

alias page="${zsb}.page nvim"
alias cpage="${zsb}.page zsb_cat"
alias cppage="${zsb}.page ccp"

hisIgnore cdpage cdnotebook page cpage cppage
