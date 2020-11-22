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
  local currentCompletion=("${COMP_WORDS[@]:1:$COMP_CWORD-1}")
  local arrayFromCommand
  local undesiredCompletion=("new" "full")

  if [ "${COMP_WORDS[1]}" = "." ]; then
    return 0
  fi

  if [ "${COMP_WORDS[1]}" = "new" ]; then
    arrayFromCommand=( $(${zsb}_getGitUntrackedFiles) )
  else
    arrayFromCommand=( $(${zsb}_getGitUnstagedFiles) )
  fi

  # Hack: by appending twice, it gets ignored by `getNonRepeatedItems`
  undesiredCompletion=("${undesiredCompletion[@]}" "${undesiredCompletion[@]}")
  local joined=("${currentCompletion[@]}" "${arrayFromCommand[@]}" "${undesiredCompletion[@]}")
  local completionArray=( $(${zsb}_getNonRepeatedItems ${joined[@]}) )

  COMPREPLY=( $(compgen -W "${completionArray[*]}") )
}




complete -F _${zsb}_ga ga
