gd() {
  if [ -z "$1" ]; then
    git diff && ${zsb}_gitStatus
    return 0
  fi

  git diff "$@" && ${zsb}_gitStatus
  return 0
}

_${zsb}_gd() {
  local currentCompletion=("${COMP_WORDS[@]:1:$COMP_CWORD-1}")
  local array=( $(${zsb}_getGitUnstagedFiles) )
  local joined=("${currentCompletion[@]}" "${array[@]}")

  local completionArray=( $(${zsb}_getNonRepeatedItems ${joined[@]}) )

  COMPREPLY=( $(compgen -W "${completionArray[*]}") )
}

complete -F _${zsb}_gd gd
