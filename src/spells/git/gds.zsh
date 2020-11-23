gds() {
  if [ -z "$1" ]; then
    git diff --staged && ${zsb}_gitStatus
    return 0
  fi

  git diff --staged "$@" && ${zsb}_gitStatus
  return 0
}

complete -C "${zsb}_getGitStagedFiles" gds

_${zsb}_gds() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList=( $(${zsb}_getGitStagedFiles) )
  local newCompletion=( $(${zsb}_removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )

  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}_gds gds
