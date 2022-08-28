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

${zsb}.getGitFiles.getgitflesfromregex() {
  ${zsb}.getGitFiles.gitShortStatus | rg "$1" | ${zsb}.getGitFiles.removeGitTokens
}

declare -gAr ZSB_FILE_TOCMD=(
  ['staged']="${zsb}.getGitFiles.getgitflesfromregex '^[MARCD]'"
  ['green']="${zsb}.getGitFiles.getgitflesfromregex '^[MARCD]'"
  ['unstaged']="${zsb}.getGitFiles.getgitflesfromregex '^.[MARCD]'"
  ['untracked']="${zsb}.getGitFiles.getgitflesfromregex '^\?{2}'"
  ['red-safe']="${zsb}.getGitFiles.getgitflesfromregex '^.[MARCD\?]'"
  ['red']="${zsb}.getGitFiles.getgitflesfromregex '^.[MARCUD\?]'"
  ['red-with-diff']="${zsb}.getGitFiles.getgitflesfromregex '^.[MARCUD]'"
  ['unmerged']="${zsb}.getGitFiles.getgitflesfromregex '(^U)|(^.U)'"
)

${zsb}.getGitFiles() (
  local this="$0"
  local newCmd="${ZSB_FILE_TOCMD[$1]}"

  if [[ -n "$newCmd" ]]
    then eval ${newCmd}
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

