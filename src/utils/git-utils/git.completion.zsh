${zsb}.gitBranches() {
  ! ${zsb}.isGitRepo && return 0
  case "$1" in
    'current')
      git branch --show-current ;;
    *)
      git branch | sed 's/^\*//' ;;
  esac
}

${zsb}.getGitFiles() {
  ! ${zsb}.isGitRepo && return 0
  case "$1" in
    'staged'|'green')
      git status --short | grep '^[MARCD]' | sed s/^...// ;;
    'unstaged')
      git status --short | grep -E '(^ [MARCD])|(^[MARCD]{2})' | sed s/^...// ;;
    'untracked')
      git status --short | grep -E '(^ \?)|(^\?{2})' | sed s/^...// ;;
    'red')
      git status --short | grep -E '(^ [MARCD\?])|(^[MARCD\?]{2})' | sed s/^...// ;;
    *)
      git status --short | sed s/^...// ;;
  esac
}
