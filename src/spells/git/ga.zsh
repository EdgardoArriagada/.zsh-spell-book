ga() (
  local this="$0"
  local firstArg="$1" # 'fast' | 'new' | '.' | @unstagedFile | @untrackedFile
  local gitFileType
  local filesToAdd=( )

  ${this}.main() {
    case "$firstArg" in
      '.') ${this}.addFilesWithDotArg ;;
      'new') ${this}.addFilesWithNewArg "$@" ;;
      'fast') ${this}.addFilesWithFastArg "$@" ;;
      *) ${this}.addFilesWithDefaulBehavior "$@" ;;
    esac
  }

  ${this}.addFilesWithDotArg() git add .

  ${this}.addFilesWithNewArg() {
    shift 1 # remove 'new' flag
    ${this}.setGitFileType 'untracked'
    ${this}.setFilesToAdd "$@"
    ${this}.addFiles
  }

  ${this}.addFilesWithFastArg() {
    shift 1 # remove 'fast' flag
    ${this}.setGitFileType 'unstaged'
    ${this}.setFilesToAdd "$@"
    ${this}.addFiles
  }

  ${this}.addFilesWithDefaulBehavior() {
    ${this}.setGitFileType 'unstaged'
    ${this}.setFilesToAdd "$@"
    ${this}.addFiles '-p'
  }

  ${this}.setGitFileType() gitFileType="$1"

  ${this}.addFiles() {
    local gitAddFlag="$1"

    ${this}.validateFilesToAdd

    [[ -z "$gitAddFlag" ]] &&
      git add "${filesToAdd[@]}" && return $?

    git add "$gitAddFlag" "${filesToAdd[@]}"
  }

  ${this}.setFilesToAdd() {
    if [[ "$#"  = "0" ]]; then
      filesToAdd=( $(${zsb}.getGitFiles "$gitFileType") )
    else
      filesToAdd=( "$@" )
    fi
  }

  ${this}.validateFilesToAdd() {
    [[ -n "$filesToAdd" ]] && return 0
    ${zsb}.cancel "There are no $(hl "$gitFileType") files to add."
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
    '.')
      return 0 ;;
    *)
      completionList=( $(${zsb}.getGitFiles 'unstaged') ) ;;
  esac

  # if we are completing the first item
  if [[ "$CURRENT" = "2" ]]; then
    case "$currentCompletion" in
      n*) # matches 'new'
        completionList+=( 'new' ) ;;
      f*) # matches 'fast'
        completionList+=( 'fast' ) ;;
    esac
  fi

  local newCompletion=( ${completionList:|usedCompletion} )

  _${zsb}.verticalComp "newCompletion"
}

compdef _${zsb}.ga ga

