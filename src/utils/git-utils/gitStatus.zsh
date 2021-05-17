${zsb}.gitStatus() {
  echo " "
  local -r gitStatusOutput=$(git -c color.status=always status --short)
  [[ -z "$gitStatusOutput" ]] && git status || echo "$gitStatusOutput"
  echo " "
}
