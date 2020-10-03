${zsb}_gitStatus() {
  echo " "

  local ${zsb}_gitStatusOutput=$(script -qc "git status --short" /dev/null < /dev/null)
  if [ -z "$${zsb}_gitStatusOutput" ]; then
    git status
  else
    echo "$${zsb}_gitStatusOutput"
  fi

  echo " "
}
