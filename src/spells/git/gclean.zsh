gclean() {
  if [ ! -z "$1" ]; then
    git clean -fd "$@"
    return $?
  fi

  echo "${ZSB_WARNING} All the added files and folders will be deleted."
  echo ""
  echo "${ZSB_PROMPT} Are you sure? [Y/n]"

  ${zsb}_yesNoMenu git clean -fd
}

complete -C "$ZSB_GIT_NEW_FILES" gclean
