gclean() {
  if [ ! -z "$1" ]; then
    git clean -fd "$@"
    return $?
  fi

  echo "${ZSB_WARNING} All the added files and folders will be deleted."
  echo ""
  echo "${ZSB_PROMPT} Are you sure? [Y/n]"

  ${zsb}.confirmMenu && git clean -fd
}

_${zsb}.gclean() {
  local usedCompletion=( "${words[@]:1:$CURRENT-2}" )
  local completionList=( $(${zsb}.getGitFiles 'untracked') )

  local newCompletion=( ${completionList:|usedCompletion} )
  _describe 'command' newCompletion
}

compdef _${zsb}.gclean gclean

