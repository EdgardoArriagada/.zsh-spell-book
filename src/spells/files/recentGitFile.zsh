${zsb}_recentGitFile() (
  local callback="$1"
  local inputFile="$2"
  eval "$callback $inputFile"
)

complete -C "$ZSB_GIT_MODIFIED_FILES" ${zsb}_recentGitFile

alias vimr="${zsb}_recentGitFile vims"
alias vr="vimr"

alias catr="${zsb}_recentGitFile cat"
alias cr="catr"
