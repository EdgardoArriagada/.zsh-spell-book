# I only use this to move staged files back to unstaged state
grh() {
  git reset -q "$@" && ${zsb}.gitStatus
}

_${zsb}.grh() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList=( $(${zsb}.getGitFiles 'staged') )
  local newCompletion=( $(${zsb}.removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )

  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}.grh grh
