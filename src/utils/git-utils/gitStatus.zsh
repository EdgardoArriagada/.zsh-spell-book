${zsb}.gitStatus() {
  <<< ''

  local gitStatusOutput=`git -c color.status=always status --short`

  if [[ -z "$gitStatusOutput" ]]
    then <<< 'nothing to commit, working tree clean'
    else <<< ${gitStatusOutput}
  fi

  <<< ''
}
