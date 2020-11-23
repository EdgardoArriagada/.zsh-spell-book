gd() {
  if [ -z "$1" ]; then
    git diff && ${zsb}_gitStatus
    return 0
  fi

  git diff "$@" && ${zsb}_gitStatus
  return 0
}

_${zsb}_gd() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList=( $(${zsb}_getGitUnstagedFiles) )
  local newCompletion=( $(${zsb}_removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )

  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}_gd gd
