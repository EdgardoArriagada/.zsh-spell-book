# I only use this to move staged files back to unstaged state
grh() {
  git reset -q "$@" && ${zsb}_gitStatus
}

_${zsb}_grh() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList=( $(${zsb}_getGitStagedFiles) )
  local newCompletion=( $(${zsb}_removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )

  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}_grh grh
