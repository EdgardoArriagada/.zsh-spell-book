${zsb}.recentGitFile() eval "$@"

_${zsb}.recentGitFile() {
  [[ "$CURRENT" -gt "3" ]] && return 0

  local completionList=( $(${zsb}.getGitFiles) )
  _${zsb}.verticalComp "completionList"
}

compdef _${zsb}.recentGitFile ${zsb}.recentGitFile

alias vr="${zsb}.recentGitFile v"
alias cr="${zsb}.recentGitFile c"
alias ccpr="${zsb}.recentGitFile ccp"
alias tigr="${zsb}.recentGitFile tig"
alias rmr="${zsb}.recentGitFile 'rm -rf'"
alias mvr="${zsb}.recentGitFile mv"
alias lr="${zsb}.recentGitFile l"

