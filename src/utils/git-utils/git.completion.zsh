${zsb}.gitBranches() {
  case $1 in
    'current') git branch --show-current 2>/dev/null ;;
    *) git branch --format='%(refname:short)' 2>/dev/null ;;
  esac
}

${zsb}.getGitFiles.removeGitTokens() sd '^...' '$1'

# '--short' is better than '--porcelain' because
# it keeps the paths of the git files relative
# to current working folder
${zsb}.getGitFiles.gitShortStatus() git status --short 2>/dev/null

${zsb}.getGitFiles.getGitFilesFromRegex() {
  ${zsb}.getGitFiles.gitShortStatus | rg "$1" | ${zsb}.getGitFiles.removeGitTokens
}

${zsb}.getGitFiles() (
  local this=$0
  local regex="${ZSB_GIT_FILETYPE_TO_REGEX[$1]}"

  if [[ -n "$regex" ]]
    then ${this}.getGitFilesFromRegex ${regex}
    else ${this}.gitShortStatus | ${this}.removeGitTokens
  fi
)

_${zsb}.gitUnrepeat() {
  # $1 can be 'staged|unstaged|untracked|red'
  local usedCompletion=( "${words[@]:1:$CURRENT-2}" )
  local completionList=( $(${zsb}.getGitFiles $1) )

  local newCompletion=( ${completionList:|usedCompletion} )

  _describe 'command' newCompletion
}

