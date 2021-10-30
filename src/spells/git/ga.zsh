ga() (
  local this="$0"
  local firstArg="$1" # 'fast' | 'new' | '.' | @unstagedFile | @untrackedFile
  local gitFileType
  local filesToAdd=( )
  local gitAddFlag

  ${this}.main() {
    case "$firstArg" in
      '.')
        ${this}.addFilesWithDotArg
        return 0
        ;;
      'new')
        shift 1 # remove 'new' flag
        ${this}.setGitFileType 'untracked'
        ${this}.setFilesToAdd "$@"
        ;;
      'fast')
        shift 1 # remove 'fast' flag
        ${this}.setGitFileType 'unstaged'
        ${this}.setFilesToAdd "$@"
        ;;
      *)
        ${this}.setGitFileType 'unstaged'
        ${this}.setFilesToAdd "$@"
        ${this}.setGitAddFlag '-p'
        ;;
    esac

    if ! ${this}.thereAreFilesToAdd; then
      ${this}.informNoChanges; return $?
    fi

    ${this}.addFiles
  }

  ${this}.addFilesWithDotArg() git add .

  ${this}.setGitFileType() gitFileType="$1"

  ${this}.addFiles() {
    if [[ -z "$gitAddFlag" ]]; then
      git add "${filesToAdd[@]}"
      return $?
    fi

    git add "$gitAddFlag" "${filesToAdd[@]}"
  }

  ${this}.setFilesToAdd() {
    if [[ "$#"  = "0" ]]; then
      filesToAdd=( $(${zsb}.getGitFiles "$gitFileType") )
    else
      filesToAdd=( "$@" )
    fi
  }

  ${this}.setGitAddFlag() gitAddFlag="$1"

  ${this}.thereAreFilesToAdd() [[ "${#filesToAdd[@]}" -gt 0 ]]

  ${this}.informNoChanges() {
    echo "${ZSB_INFO} There are no $(hl "$gitFileType") files to add."
    return 0
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
  local formattedComp=( $(${zsb}.formatComp "newCompletion") )
  _describe 'command' formattedComp
}

compdef _${zsb}.ga ga

