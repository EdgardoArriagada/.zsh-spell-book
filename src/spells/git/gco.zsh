# I only use this for checking out files
gco() {
  git checkout "$@" && ${zsb}.gitStatus
}

_${zsb}.gco() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList=( $(${zsb}.getGitUnstagedFiles) )
  local newCompletion=( $(${zsb}.removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )

  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}.gco gco
