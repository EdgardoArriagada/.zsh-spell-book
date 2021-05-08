gds() {
  if [ -z "$1" ]; then
    git diff --staged && ${zsb}.gitStatus
    return 0
  fi

  git diff --staged "$@" && ${zsb}.gitStatus
  return 0
}

complete -C "${zsb}.getGitFiles 'staged'" gds

_${zsb}.gds() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList=( $(${zsb}.getGitFiles 'staged') )
  local newCompletion=( $(${zsb}.removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )

  COMPREPLY=( "${newCompletion[@]}" )
}

complete -F _${zsb}.gds gds
