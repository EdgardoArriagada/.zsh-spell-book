${zsb}.gitStatus() {
  print ''

  local gitStatusOutput=`git -c color.status=always status --short`

  if [[ -z "$gitStatusOutput" ]] 
    then print 'nothing to commit, working tree clean'
    else print ${gitStatusOutput}
  fi

  print ''
}
