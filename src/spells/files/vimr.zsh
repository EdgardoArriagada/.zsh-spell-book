# vim resent modified files in git
vimr() {
  vim "$@" && isGitRepo && gitStatus
}

# Complete with git status files relative to current directory
complete -C "isGitRepo && git status --short | sed s/^...// || ls" vimr

alias vir="vimr"
alias vr="vimr"
