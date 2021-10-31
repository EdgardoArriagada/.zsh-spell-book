${zsb}.gitBranches() {
  ${zsb}.isGitRepo || return 0
  case "$1" in
    'current')
      git branch --show-current ;;
    *)
      git branch | sed 's/^\*//' ;;
  esac
}

${zsb}.getGitFiles() (
  local this="$0"

  ${this}.main() {
    ${zsb}.isGitRepo || return 0
    case "$1" in
      'staged'|'green')
        ${this}.getgitflesfromregex '^[MARCD]' ;;
      'unstaged')
        ${this}.getgitflesfromregex '^.[MARCD]' ;;
      'untracked')
        ${this}.getgitflesfromregex '^\?{2}' ;;
      'red-safe')
        ${this}.getgitflesfromregex '^.[MARCD\?]' ;;
      'red')
        ${this}.getgitflesfromregex '^.[MARCUD\?]' ;;
      'red-with-diff')
        ${this}.getgitflesfromregex '^.[MARCUD]' ;;
      'unmerged')
        ${this}.getgitflesfromregex '(^U)|(^.U)' ;;
      *)
        ${this}.gitShortStatus | ${this}.removeGitTokens ;;
    esac
  }

  ${this}.getgitflesfromregex() {
    ${this}.gitShortStatus | rg "$1" | ${this}.removeGitTokens
  }

  # '--short' is better than '--porcelain' because
  # it keeps the paths of the git files relative
  # to current working folder
  ${this}.gitShortStatus() git status --short

  ${this}.removeGitTokens() sed 's/^...//'

  ${this}.main "$@"
)

_${zsb}.gitUnrepeat() {
  # $1 can be 'staged|unstaged|untracked|red'
  local usedCompletion=( "${words[@]:1:$CURRENT-2}" )
  local completionList=( $(${zsb}.getGitFiles $1) )

  local newCompletion=( ${completionList:|usedCompletion} )

  _${zsb}.verticalComp "newCompletion"
}

