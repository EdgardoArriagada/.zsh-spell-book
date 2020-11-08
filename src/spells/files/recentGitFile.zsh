${zsb}_recentGitFile() (
  local callback="$1"
  local inputFile="$2"
  eval "$callback $inputFile"
)

# Complete with git status files relative to current directory
complete -C "${zsb}_isGitRepo && git status --short | sed s/^...// || ls" ${zsb}_recentGitFile

alias vimr="${zsb}_recentGitFile vims"
alias vr="vimr"

alias catr="${zsb}_recentGitFile cat"
alias cr="catr"
