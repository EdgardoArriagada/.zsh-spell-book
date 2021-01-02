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

  ${this}.gitShortStatus() git status --porcelain

  ${this}.removeGitTokens() sed 's/^...//'

  ${this}.main "$@"
)
