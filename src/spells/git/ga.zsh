ga() (
  local this="$0"
  local firstArg="$1" # 'fast' | 'new' | '.' | @unstagedFile | @untrackedFile
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
    ${this}.setFilesToAdd 'untracked' "$@"
    ${this}.addFiles
  }

  ${this}.addFilesWithFastArg() {
    shift 1 # remove 'fast' flag
    ${this}.setFilesToAdd 'unstaged' "$@"
    ${this}.addFiles
  }

  ${this}.addFilesWithDefaulBehavior() {
    ${this}.setFilesToAdd 'unstaged' "$@"
    ${this}.addFiles '-p'
  }

  ${this}.setGitFileType() gitFileType="$1"

  ${this}.addFiles() {
    if [[ -n "$1" ]]; then
      git add "$@" "${filesToAdd[@]}"; return $?
    fi

    git add "${filesToAdd[@]}"
  }

  ${this}.setFilesToAdd() {
    local gitFileType="$1"; shift 1

    if [[ "$#" = "0" ]]; then
      filesToAdd=( $(${zsb}.getGitFiles "$gitFileType") )
    else
      filesToAdd=( "$@" )
    fi

    ${this}.validateFilesToAdd
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

