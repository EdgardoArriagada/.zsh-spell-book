ga() {
  if [ "$1" = "." ]; then
    git add . && ${zsb}_gitStatus
    return 0
  fi

  if [ "$1" = "new" ]; then
    if [ -z "$2" ]; then
      echo "${ZSB_ERROR} Untracked file(s) expected."
      return 1
    fi
    shift 1 # remove 'new' flag
    git add "$@" && ${zsb}_gitStatus
    return 0
  fi

  if [ "$1" = "full" ]; then
    if [ -z "$2" ]; then
      echo "${ZSB_ERROR} Unstaged file(s) expected."
      return 1
    fi
    shift 1 # remove 'full' flag
    git add "$@" && ${zsb}_gitStatus
    return 0
  fi

  git add -p "$@" && ${zsb}_gitStatus
  return 0
}

_${zsb}_ga() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList

  if [ "${COMP_WORDS[1]}" = "." ]; then
    return 0
  fi

  if [ "${COMP_WORDS[1]}" = "new" ]; then
    completionList=( $(${zsb}_getGitUntrackedFiles) )
  else
    completionList=( $(${zsb}_getGitUnstagedFiles) )
  fi

  local newCompletion=( $(${zsb}_removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )
  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}_ga ga
