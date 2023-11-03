${zsb}.gitBranches() {
  case $1 in
    'current') git branch --show-current 2>/dev/null ;;
    *) git branch --format='%(refname:short)' 2>/dev/null ;;
  esac
}

${zsb}.getGitFiles() (
  # '--short' is better than '--porcelain' because
  # it keeps the paths of the git files relative
  # to current working folder

  case $# in
    0) git status --short 2>/dev/null | cut -c 4- ;;
    *) git status --short 2>/dev/null | rg ${ZSB_GIT_FILETYPE_TO_REGEX[$1]} | cut -c 4- ;;
  esac
)

_${zsb}.gitUnrepeat() {
  # $1 can be 'staged|unstaged|untracked|red'
  local usedCompletion=( ${words[@]:1:$CURRENT-2} )
  local completionList=( $(${zsb}.getGitFiles $1) )

  local newCompletion=( ${completionList:|usedCompletion} )

  _describe 'command' newCompletion
}

