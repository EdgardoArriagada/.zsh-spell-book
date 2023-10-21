export ZSB_CHAPT=${ZSB_NOTEBOOK_DIR}/${ZSB_NOTEBOOK_CHAPTER}

${zsb}.page.validateEnv() {
  ${zsb}.validate "ZSB_NOTEBOOK_DIR"
  ${zsb}.validate "ZSB_NOTEBOOK_CHAPTER"
}

${zsb}.page.getTodaysPage() {
  local year=`date +%Y`
  local month=`date +%b`
  local page=`date +%b-%d-%Y`

  local currentPage=${ZSB_CHAPT}/${year}/${month}/${page}
  mkdir -p `dirname $currentPage`

  <<< $currentPage
}


${zsb}.page() {
  local cmd=$1
  local currentPage=$2
  ${zsb}.page.validateEnv

  if [[ -z "$currentPage" ]]; then
    local todaysPage=`${zsb}.page.getTodaysPage`

    (builtin cd `dirname $todaysPage` && \
      $cmd ${todaysPage}.md)
  else
    (builtin cd $ZSB_CHAPT && \
      $cmd ${currentPage}.md)
  fi
}

_${zsb}.page() {
  local files=(${(@f)$(ls $ZSB_CHAPT | rg --invert-match '__off\.md$')})
  local result=()

  for f in $files; do
    if [[ -f ${ZSB_CHAPT}/${f} ]];then
      # trim the extension
       result+=(${f%.*})
    fi
  done

  _describe 'command' result
}

compdef _${zsb}.page ${zsb}.page

alias cdnotebook="${zsb}.validate 'ZSB_NOTEBOOK_DIR' && cds ${ZSB_NOTEBOOK_DIR}"
alias page="${zsb}.page nvim"
alias cpage="${zsb}.page zsb_cat"
alias ccpage="${zsb}.page ccp"
alias cdpage="${zsb}.page.validateEnv && cd ${ZSB_CHAPT}"

hisIgnore cdnotebook page cpage cppage cdpage
