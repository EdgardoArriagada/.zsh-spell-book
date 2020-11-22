gpristine(){
  echo "${ZSB_WARNING} All your uncommited work will be lost"
  echo ""
  echo "${ZSB_PROMPT} Are you sure? [Y/n]"
  ${zsb}_yesNoMenu && \
    git reset --hard && git clean -dffx
}

complete gpristine
