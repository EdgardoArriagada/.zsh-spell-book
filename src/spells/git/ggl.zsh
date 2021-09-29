ggl() {
  local -r inputBranch=${1:?'You must provide a branch to pull from.'}
  local -r currentBranch=$(${zsb}.gitBranches 'current')

  if [[ "$currentBranch" != "$inputBranch" ]]; then
    ${zsb}.throw "You can't use this command to pull from different branches"
  fi

  git pull origin "$inputBranch"
}

compdef "_${zsb}.nonRepeatedListC 'git branch --show-current'" ggl

