gd() {
  if [ -z "$1" ]; then
    git diff && ${zsb}.gitStatus
    return 0
  fi

  git diff "$@" && ${zsb}.gitStatus
  return 0
}

_${zsb}.gd() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList=( $(${zsb}.getGitUnstagedFiles) )
  local newCompletion=( $(${zsb}.removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )

  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}.gd gd
