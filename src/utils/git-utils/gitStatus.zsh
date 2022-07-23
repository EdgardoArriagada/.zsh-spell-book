${zsb}.gitStatus() {
  print " "
  local -r gitStatusOutput=$(git -c color.status=always status --short)
  [[ -z "$gitStatusOutput" ]] && print "nothing to commit, working tree clean" || print "$gitStatusOutput"
  print " "
}
