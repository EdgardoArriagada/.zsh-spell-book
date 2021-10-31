ga() (
  local this="$0"
  local firstArg="$1" # 'fast' | 'new' | 'unmerged' | '.' | @unstagedFile | @untrackedFile
  local filesToAdd=( )

  ${this}.main() {
    case "$firstArg" in
      '.') ${this}.addFilesWithDotFlag ;;
      'new') ${this}.addFilesWithNewFlag "$@" ;;
      'fast') ${this}.addFilesWithFastFlag "$@" ;;
      'unmerged') ${this}.addFilesWithUnmergedFlag "$@" ;;
      *) ${this}.addFilesWithDefaulBehavior "$@" ;;
    esac
  }

  ${this}.addFilesWithDotFlag() git add .

  ${this}.addFilesWithNewFlag() {
    shift 1 # remove 'new' flag
    ${this}.setFilesToAdd 'untracked' "$@"
    ${this}.addFiles
  }

  ${this}.addFilesWithFastFlag() {
    shift 1 # remove 'fast' flag
    ${this}.setFilesToAdd 'unstaged' "$@"
    ${this}.addFiles
  }

  ${this}.addFilesWithUnmergedFlag() {
    shift 1 # remove 'unmerged' flag
    ${this}.setFilesToAdd 'unmerged' "$@"
    ${this}.addFiles
  }

  ${this}.addFilesWithDefaulBehavior() {
    (( $# = 0 )) && galock
    ${this}.setFilesToAdd 'unstaged' "$@"
    ${this}.addFiles '-p'
  }

  ${this}.setFilesToAdd() {
    local gitFileType="$1"; shift 1

    if (( $# = 0 )); then
      filesToAdd=( $(${zsb}.getGitFiles "$gitFileType") )
    else
      filesToAdd=( "$@" )
    fi

    ${this}.validateFilesToAdd "$gitFileType"
  }

  ${this}.validateFilesToAdd() {
    [[ -z "$filesToAdd" ]] &&
      ${zsb}.cancel "There are no $(hl $1) files to add."
  }

  ${this}.addFiles() {
    if [[ -n "$1" ]]; then
      git add "$@" "${filesToAdd[@]}"; return $?
    fi

    git add "${filesToAdd[@]}"
  }

  ${this}.main "$@"
  ${zsb}.gitStatus
)

_${zsb}.ga() {
  local usedCompletion=( "${words[@]:1:$CURRENT-2}" )
  local firstItemUsed="${words[2]}"
  local currentCompletion="${words[CURRENT]}"
  local completionList=( )

  case "$firstItemUsed" in
    'new')
      completionList=( $(${zsb}.getGitFiles 'untracked') ) ;;
    'fast')
      completionList=( $(${zsb}.getGitFiles 'unstaged') ) ;;
    'unmerged')
      completionList=( $(${zsb}.getGitFiles 'unmerged') ) ;;
    '.')
      return 0 ;;
    *)
      completionList=( $(${zsb}.getGitFiles 'unstaged') ) ;;
  esac

  # if we are completing the first item
  if [[ "$CURRENT" = "2" ]]; then
    case "$currentCompletion" in
      n*) completionList+=('new') ;;
      f*) completionList+=('fast') ;;
      u*) completionList+=('unmerged') ;;
    esac
  fi

  local newCompletion=( ${completionList:|usedCompletion} )

  _${zsb}.verticalComp "newCompletion"
}

compdef _${zsb}.ga ga

