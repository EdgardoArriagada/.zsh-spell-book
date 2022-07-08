${zsb}.gitStatus() {
  print " "
  local -r gitStatusOutput=$(git -c color.status=always status --short)
  [[ -z "$gitStatusOutput" ]] && git status || print "$gitStatusOutput"
  print " "
}
