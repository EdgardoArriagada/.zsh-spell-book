${zsb}_gitStatus() {
  echo " "

  local gitStatusOutput=$(script -qc "git status --short" /dev/null < /dev/null)
  if [ -z "$gitStatusOutput" ]; then
    git status
  else
    echo "$gitStatusOutput"
  fi

  echo " "
}
