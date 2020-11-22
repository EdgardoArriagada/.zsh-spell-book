# I only use this for checking out files
gco() {
  git checkout "$@" && ${zsb}_gitStatus
}

_${zsb}_gco() {
  local currentCompletion=("${COMP_WORDS[@]:1:$COMP_CWORD-1}")
  local array=( $(${zsb}_getGitUnstagedFiles) )
  local joined=("${currentCompletion[@]}" "${array[@]}")

  local completionArray=( $(${zsb}_getNonRepeatedItems ${joined[@]}) )

  COMPREPLY=( $(compgen -W "${completionArray[*]}") )
}

complete -F _${zsb}_gco gco
