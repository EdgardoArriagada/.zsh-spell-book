gds() {
  if [ -z "$1" ]; then
    git diff --staged && ${zsb}.gitStatus
    return 0
  fi

  git diff --staged "$@" && ${zsb}.gitStatus
  return 0
}

complete -C "${zsb}.getGitStagedFiles" gds

_${zsb}.gds() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList=( $(${zsb}.getGitStagedFiles) )
  local newCompletion=( $(${zsb}.removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )

  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}.gds gds
