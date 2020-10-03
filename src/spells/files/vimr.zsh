# vim resent modified files in git
vimr() {
  vim "$@" && ${zsb}_isGitRepo && ${zsb}_gitStatus
}

# Complete with git status files relative to current directory
complete -C "${zsb}_isGitRepo && git status --short | sed s/^...// || ls" vimr

alias vir="vimr"
alias vr="vimr"
