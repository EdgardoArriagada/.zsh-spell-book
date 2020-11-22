# I only use this to move staged files back to unstaged state
grh() {
  git reset -q "$@" && ${zsb}_gitStatus
}

_${zsb}_grh() {
  local currentCompletion=("${COMP_WORDS[@]:1:$COMP_CWORD-1}")
  local array=( $(${zsb}_getGitStagedFiles) )
  local joined=("${currentCompletion[@]}" "${array[@]}")

  local completionArray=( $(${zsb}_getNonRepeatedItems ${joined[@]}) )

  COMPREPLY=( $(compgen -W "${completionArray[*]}") )
}

complete -F _${zsb}_grh grh
