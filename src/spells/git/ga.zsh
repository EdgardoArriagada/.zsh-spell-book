ga() (
  local this="$0"
  local firstArg="$1" # 'fast' | 'new' | '.' | @unstagedFile | @untrackedFile
  local filesToAdd=( )

  ${this}.main() {
    case "$firstArg" in
      '.')
        ${this}.addFilesWithDotArg ;;
      'new')
        ${this}.addFilesWithNewArg "$@" ;;
      'fast')
        ${this}.addFilesWithFastArg "$@" ;;
      *)
        ${this}.addFilesDefault "$@" ;;
    esac
  }

  ${this}.addFilesWithDotArg() git add .

  ${this}.addFilesWithNewArg() {
    shift 1 # remove 'new' flag
    if [ "$#"  = "0" ]; then
      filesToAdd=( $(${zsb}.getGitFiles 'untracked') )
    else
      filesToAdd=( "$@" )
    fi

    if ! ${this}.thereAreFilesToAdd; then
      ${this}.informNoChanges 'untracked'; return $?
    fi

    git add "${filesToAdd[@]}"
  }

  ${this}.addFilesWithFastArg() {
    shift 1 # remove 'fast' flag
    if [ "$#"  = "0" ]; then
      filesToAdd=( $(${zsb}.getGitFiles 'unstaged') )
    else
      filesToAdd=( "$@" )
    fi

    if ! ${this}.thereAreFilesToAdd; then
      ${this}.informNoChanges 'unstaged'; return $?
    fi

    git add "${filesToAdd[@]}"
  }

  ${this}.thereAreFilesToAdd() [ "${#filesToAdd[@]}" -gt 0 ]

  ${this}.informNoChanges() {
    local gitFileType="$1"
    echo "${ZSB_INFO} There are no $(hl "$gitFileType") files to add."
    return 0
  }

  ${this}.addFilesDefault() {
    if [ "$#"  = "0" ]; then
      filesToAdd=( $(${zsb}.getGitFiles 'unstaged') )
    else
      filesToAdd=( "$@" )
    fi

    if ! ${this}.thereAreFilesToAdd; then
      ${this}.informNoChanges 'unstaged'; return $?
    fi

    git add -p "${filesToAdd[@]}"
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
      completionList=( $(${zsb}.getGitFiles 'untracked') ) ;;
    'fast')
      completionList=( $(${zsb}.getGitFiles 'unstaged') ) ;;
    '.')
      return 0 ;;
    *)
      completionList=( $(${zsb}.getGitFiles 'unstaged') ) ;;
  esac

  # if we are completing the first item
  if [ "$COMP_CWORD" = "1" ]; then
    case "$currentCompletion" in
      [n]*) # matches 'new'
        completionList+=( 'new' ) ;;
      [f]*) # matches 'fast'
        completionList=( 'fast' ) ;;
    esac
  fi

  local newCompletion=( $(${zsb}.removeUsedOptions "${usedCompletion[*]}" "${completionList[*]}") )
  COMPREPLY=( $(compgen -W "${newCompletion[*]}") )
}

complete -F _${zsb}.ga ga
