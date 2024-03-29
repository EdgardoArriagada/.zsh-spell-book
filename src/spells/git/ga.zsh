# $1 can be 'fast' | 'new' | 'unmerged' | '.' | @unstagedFile | @untrackedFile
ga() (
  local this=$0
  local -a filesToAdd

  ${this}.addFilesWithDotFlag() git add .

  ${this}.addFilesWithNewFlag() {
    shift 1 # remove 'new' flag
    ${this}.setFilesToAdd 'untracked' $@
    (( $# > 0 )) && ${this}.areAddingPackageFiles && galock
    ${this}.addFiles
  }

  ${this}.addFilesWithFastFlag() {
    shift 1 # remove 'fast' flag
    ${this}.setFilesToAdd 'unstaged' $@
    ${this}.areAddingPackageFiles && galock
    ${this}.addFiles
  }

  ${this}.addFilesWithUnmergedFlag() {
    shift 1 # remove 'unmerged' flag
    ${this}.setFilesToAdd 'unmerged' $@
    ${this}.addFiles
  }

  ${this}.addFilesWithDefaulBehavior() {
    ${this}.setFilesToAdd 'unstaged' $@
    ${this}.areAddingPackageFiles && galock
    ${this}.addFiles '-p'
  }

  ${this}.areAddingPackageFiles() [[ -n ${ZSB_GIT_PACKAGE_FILES:*filesToAdd} ]]

  ${this}.setFilesToAdd() {
    local gitFileType=$1
    shift 1

    if (( $# ))
      then filesToAdd=( $@ )
      else filesToAdd=( $(${zsb}.getGitFiles "$gitFileType") )
    fi

    ${this}.validateFilesToAdd "$gitFileType"
  }

  ${this}.validateFilesToAdd() {
    [[ -z "$filesToAdd" ]] &&
      ${zsb}.cancel "There are no $(hl $1) files to add."
  }

  ${this}.addFiles() {
    if (( $# ))
      then git add $@ ${filesToAdd[@]}
      else git add ${filesToAdd[@]}
    fi
  }

  { # main
    case $1 in
      .) ${this}.addFilesWithDotFlag ;;
      new) ${this}.addFilesWithNewFlag "$@" ;;
      fast) ${this}.addFilesWithFastFlag "$@" ;;
      --unmerged) ${this}.addFilesWithUnmergedFlag "$@" ;;
      *) ${this}.addFilesWithDefaulBehavior "$@" ;;
    esac

    ${zsb}.gitStatus
  }
)

_${zsb}.ga() {
  local usedCompletion=(${words[@]:1:$CURRENT - 2})
  local firstItemUsed=${words[2]}
  local -a completionList

  case $firstItemUsed in
    .) return 0 ;;
    new) completionList=( $(${zsb}.getGitFiles 'untracked') ) ;;
    fast) completionList=( $(${zsb}.getGitFiles 'unstaged') ) ;;
    --unmerged) completionList=( $(${zsb}.getGitFiles 'unmerged') ) ;;
    *) completionList=( $(${zsb}.getGitFiles 'unstaged') ) ;;
  esac

  # if we are completing the first item
  if (( $CURRENT == 2 )) then
    case $firstItemUsed in
      n*) completionList+=(new) ;;
      f*) completionList+=(fast) ;;
      -*) completionList+=(--unmerged) ;;
    esac
  fi

  local newCompletion=(${completionList:|usedCompletion})

  _describe 'command' newCompletion
}

hisIgnore ga 'ga .'

compdef _${zsb}.ga ga

