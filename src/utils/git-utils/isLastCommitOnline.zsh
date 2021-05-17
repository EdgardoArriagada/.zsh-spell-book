${zsb}.isLastCommitOnline() {
  local WORKING_BRANCH=$(git symbolic-ref --short HEAD)
  local HEAD_HASH=$(git rev-parse HEAD)
  local REMOTE_HEAD_HASH=$(git rev-parse --revs-only origin/${WORKING_BRANCH})

  [[ "$HEAD_HASH" = "$REMOTE_HEAD_HASH" ]]
}
