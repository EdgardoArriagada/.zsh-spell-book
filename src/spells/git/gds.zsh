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
  local currentCompletion=("${COMP_WORDS[@]:1:$COMP_CWORD-1}")
  local array=( $(${zsb}_getGitStagedFiles) )
  local joined=("${currentCompletion[@]}" "${array[@]}")

  local completionArray=( $(${zsb}_getNonRepeatedItems ${joined[@]}) )

  COMPREPLY=( $(compgen -W "${completionArray[*]}") )
}

complete -F _${zsb}_gds gds
