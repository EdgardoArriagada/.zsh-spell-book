ga() (
  local firstArg="$1" # 'fast' | 'new' | '.' | @unstagedFile | @untrackedFile
  local amountOfArgs="$#"
  local filesToAdd=( )

  main() {
    if [ "$firstArg" = "." ]; then
      git add .
      return 0
    fi

    if existFlag && hasManyArgs; then
      shift 1
      git add "$@"
      return 0
    fi

    setFilesToAdd

    if ! thereAreFilesToAdd; then
      informNoChanges; return 0
    fi

    addFilesToGitStaging

    return 0
  }

  existFlag() ([ "$firstArg" = "new" ] || [ "$firstArg" = "fast" ])

  hasManyArgs() [ "$amountOfArgs" -gt 1 ]

  setFilesToAdd() {
    if [ "$firstArg" = "new" ]; then
      filesToAdd=( $(${zsb}_getGitUntrackedFiles) )
      return 0
    fi

    filesToAdd=( $(${zsb}_getGitUnstagedFiles) )
  }

  thereAreFilesToAdd() [ "${#filesToAdd[@]}" -gt 0 ]

  informNoChanges() {
    echo "${ZSB_INFO} No changes."
    return 0
  }

  addFilesToGitStaging() {
    if existFlag; then
      git add "${filesToAdd[@]}"
    else
      git add -p "${filesToAdd[@]}"
    fi
  }

  main "$@"
  ${zsb}_gitStatus
)

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
