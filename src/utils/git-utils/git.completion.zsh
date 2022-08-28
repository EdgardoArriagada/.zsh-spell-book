${zsb}.gitBranches() {
  ${zsb}.isGitRepo || return 0
  case "$1" in
    'current') git branch --show-current ;;
    *) git branch | sd '^\*' '$1' ;;
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

declare -gAr ZSB_GIT_FILETYPE_TO_REGEX=(
  ['staged']='^[MARCD]'
  ['unstaged']='^.[MARCD]'
  ['untracked']='^\?{2}'
  ['red-safe']='^.[MARCD\?]'
  ['red']='^.[MARCUD\?]'
  ['red-with-diff']='^.[MARCUD]'
  ['unmerged']='(^U)|(^.U)'
)

${zsb}.getGitFiles() (
  local this="$0"
  local token="${ZSB_GIT_FILETYPE_TO_REGEX[$1]}"

  if [[ -n "$token" ]]
    then ${this}.getGitFilesFromRegex ${token}
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

