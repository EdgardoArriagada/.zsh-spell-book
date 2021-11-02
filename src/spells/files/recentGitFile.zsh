${zsb}.recentGitFile() {
  local callback="$1"
  local inputFile="$2"
  eval "${callback} ${inputFile}"
}

_${zsb}.recentGitFile() {
  [[ "$CURRENT" -gt "3" ]] && return 0

  local completionList=( $(${zsb}.getGitFiles) )
  _${zsb}.verticalComp "completionList"
}

compdef _${zsb}.recentGitFile ${zsb}.recentGitFile

alias vr="${zsb}.recentGitFile vims"
alias cr="${zsb}.recentGitFile c"
alias ccpr="${zsb}.recentGitFile ccp"
alias tigr="${zsb}.recentGitFile tig"
alias rmr="${zsb}.recentGitFile 'rm -rf'"

