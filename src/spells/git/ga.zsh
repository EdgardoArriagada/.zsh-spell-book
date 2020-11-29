ga() (
  local this="$0"
  local firstArg="$1" # 'fast' | 'new' | '.' | @unstagedFile | @untrackedFile
  local filesToAdd=( )

  ${this}.main() {
    case "$firstArg" in
      'new')
        ${this}.addFilesWithNewArg "$@"
      ;;
      'fast')
        ${this}.addFilesWithFastArg "$@"
      ;;
      '.')
        ${this}.addFilesWithDotArg
      ;;
      *)
        ${this}.addFilesDefault
      ;;
    esac
  }

  ${this}.addFilesWithNewArg() {
    shift 1 # remove 'new' flag
    if [ "$#"  = "0" ]; then
      filesToAdd=( $(${zsb}.getGitUntrackedFiles) )
    else
      filesToAdd=( "$@" )
    fi

    if ! ${this}.thereAreFilesToAdd; then
      ${this}.informNoChanges; return $?
    fi

    git add "${filesToAdd[@]}"
  }

  ${this}.addFilesWithFastArg() {
    shift 1 # remove 'fast' flag
    if [ "$#"  = "0" ]; then
      filesToAdd=( $(${zsb}.getGitUnstagedFiles) )
    else
      filesToAdd=( "$@" )
    fi

    if ! ${this}.thereAreFilesToAdd; then
      ${this}.informNoChanges; return $?
    fi

    git add "${filesToAdd[@]}"
  }

  ${this}.addFilesWithDotArg() {
    git add .
    return 0
  }

  ${this}.addFilesDefault() {
    filesToAdd=( "$@" )

    git add -p "${filesToAdd[@]}"
  }

  ${this}.thereAreFilesToAdd() [ "${#filesToAdd[@]}" -gt 0 ]

  ${this}.informNoChanges() {
    echo "${ZSB_INFO} No changes."
    return 0
  }

  ${this}.main "$@"
  ${zsb}.gitStatus
)

_${zsb}.ga() {
  local usedCompletion=( "${COMP_WORDS[@]:1:$COMP_CWORD-1}" )
  local firstItemUsed="${COMP_WORDS[1]}"
  local currentCompletion="${COMP_WORDS[COMP_CWORD]}"
  local completionList=( )

  case "$firstItemUsed" in
    'new')
      completionList=( $(${zsb}.getGitUntrackedFiles) )
    ;;
    'fast')
      completionList=( $(${zsb}.getGitUnstagedFiles) )
    ;;
    '.')
      return 0
    ;;
    *)
      completionList=( $(${zsb}.getGitUnstagedFiles) )
    ;;
  esac

  # if we are completing the first item
  if [ "$COMP_CWORD" = "1" ]; then
    case "$currentCompletion" in
      [n]*) # matches 'new'
        completionList+=( 'new' )
      ;;
      [f]*) # matches 'fast'
        completionList=( 'fast' )
      ;;
    esac
  fi

  local newCompletion=( $(${zsb}.removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )
  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}.ga ga
