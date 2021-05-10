# I only use this to move staged files back to unstaged state
grh() {
  git reset -q "$@" && ${zsb}.gitStatus
}

_${zsb}.grh() {
  local usedCompletion=( "${words[@]:1:$CURRENT-1}" )
  local completionList=( $(${zsb}.getGitFiles 'staged') )

  local newCompletion=( ${completionList:|usedCompletion} )
  _describe 'command' newCompletion
}

compdef _${zsb}.grh grh

