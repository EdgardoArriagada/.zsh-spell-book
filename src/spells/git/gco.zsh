# I only use this for checking out files
gco() {
  git checkout "$@" && ${zsb}.gitStatus
}

_${zsb}.gco() {
  local usedCompletion=( "${words[@]:1:$CURRENT-2}" )
  local completionList=( $(${zsb}.getGitFiles 'unstaged') )

  local newCompletion=( ${completionList:|usedCompletion} )
  _describe 'command' newCompletion
}

compdef _${zsb}.gco gco

