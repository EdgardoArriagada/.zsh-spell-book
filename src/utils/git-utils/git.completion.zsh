${zsb}.gitBranches() {
  ! ${zsb}.isGitRepo && return 0
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
    ! ${zsb}.isGitRepo && return 0
    case "$1" in
      'staged'|'green')
        ${this}.getgitflesfromregex '^[MARCD]' ;;
      'unstaged')
        ${this}.getgitflesfromregex '(^ [MARCD])|(^[MARCD]{2})' ;;
      'untracked')
        ${this}.getgitflesfromregex '(^ \?)|(^\?{2})' ;;
      'red')
        ${this}.getgitflesfromregex '(^ [MARCD\?])|(^[MARCD\?]{2})' ;;
      *)
        ${this}.gitShortStatus | ${this}.removeGitTokens ;;
    esac
  }

  ${this}.getgitflesfromregex() {
    ${this}.gitShortStatus | grep -E "$1" | ${this}.removeGitTokens
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
  _describe 'command' newCompletion
}

