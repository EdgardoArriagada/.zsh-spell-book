${zsb}.gitBranches() {
  ${zsb}.isGitRepo && git branch | sed 's/^\*//'
}

${zsb}.getGitFiles() {
  ! ${zsb}.isGitRepo && return 0
  local gitFileType="$1"
  case "$gitFileType" in
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
