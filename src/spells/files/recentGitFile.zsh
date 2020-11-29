${zsb}.recentGitFile() (
  local callback="$1"
  local inputFile="$2"
  eval "$callback $inputFile"
)

_${zsb}.recentGitFile() {
  [ "$COMP_CWORD" -gt "2" ] && return 0

  COMPREPLY=( $(compgen -C "${zsb}.getGitModifiedFiles") )
}

complete -F _${zsb}.recentGitFile ${zsb}.recentGitFile

alias vimr="${zsb}.recentGitFile vims"
alias vr="vimr"

alias catr="${zsb}.recentGitFile cat"
alias cr="catr"
