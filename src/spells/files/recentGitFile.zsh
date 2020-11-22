${zsb}_recentGitFile() (
  local callback="$1"
  local inputFile="$2"
  eval "$callback $inputFile"
)

_${zsb}_recentGitFile() {
  [ "$COMP_CWORD" -gt "2" ] && return 0

  COMPREPLY=( $(compgen -C "${zsb}_getGitModifiedFiles") )
}

complete -F _${zsb}_recentGitFile ${zsb}_recentGitFile

alias vimr="${zsb}_recentGitFile vims"
alias vr="vimr"

alias catr="${zsb}_recentGitFile cat"
alias cr="catr"
