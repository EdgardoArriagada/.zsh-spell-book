# I only use this for checking out files
gco() {
  git checkout "$@" && ${zsb}_gitStatus
}

_${zsb}_gco() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList=( $(${zsb}_getGitUnstagedFiles) )
  local newCompletion=( $(${zsb}_removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )

  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}_gco gco
