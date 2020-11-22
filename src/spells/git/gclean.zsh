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

_${zsb}_gclean() {
  local currentCompletion=("${COMP_WORDS[@]:1:$COMP_CWORD-1}")
  local array=( $(${zsb}_getGitUntrackedFiles) )
  local joined=("${currentCompletion[@]}" "${array[@]}")

  local completionArray=( $(${zsb}_getNonRepeatedItems ${joined[@]}) )

  COMPREPLY=( $(compgen -W "${completionArray[*]}") )
}

complete -F _${zsb}_gclean gclean
