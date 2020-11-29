ga() (
  local this="$0"
  local firstArg="$1" # 'fast' | 'new' | '.' | @unstagedFile | @untrackedFile
  local filesToAdd=( )

  ${this}.main() {
    ${this}.tryToAddFilesWithDotArg && return 0

    ${this}.existFlag && shift 1

    if [ "$#" = "0" ]; then
      ${this}.setFilesToAddFromGitList
    else
      filesToAdd=( "$@" )
    fi

    if ! ${this}.thereAreFilesToAdd; then
      ${this}.informNoChanges; return 0
    fi

    ${this}.addFilesToGitStaging
  }

  ${this}.tryToAddFilesWithDotArg() {
    if [ "$firstArg" = "." ]; then
      git add .
      return 0
    fi
    return 1
  }

  ${this}.existFlag() ([ "$firstArg" = "new" ] || [ "$firstArg" = "fast" ])

  ${this}.setFilesToAddFromGitList() {
    if [ "$firstArg" = "new" ]; then
      filesToAdd=( $(${zsb}.getGitUntrackedFiles) )
      return 0
    fi
    filesToAdd=( $(${zsb}.getGitUnstagedFiles) )
  }

  ${this}.thereAreFilesToAdd() [ "${#filesToAdd[@]}" -gt 0 ]

  ${this}.informNoChanges() {
    echo "${ZSB_INFO} No changes."
    return 0
  }

  ${this}.addFilesToGitStaging() {
    if ${this}.existFlag; then
      git add "${filesToAdd[@]}"
    else
      git add -p "${filesToAdd[@]}"
    fi
  }

  ${this}.main "$@"
  ${zsb}.gitStatus
)

_${zsb}.ga() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local completionList

  if [ "${COMP_WORDS[1]}" = "." ]; then
    return 0
  fi

  if [ "${COMP_WORDS[1]}" = "new" ]; then
    completionList=( $(${zsb}.getGitUntrackedFiles) )
  else
    completionList=( $(${zsb}.getGitUnstagedFiles) )
  fi

  local newCompletion=( $(${zsb}.removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )
  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}.ga ga
