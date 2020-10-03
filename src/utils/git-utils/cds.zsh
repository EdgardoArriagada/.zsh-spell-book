cds() {
  echo " "
  builtin cd $1 && ls

  if ${zsb}_isGitRepo; then
    printf %"$COLUMNS"s | tr " " "_"
    ${zsb}_gitStatus
    return 0
  fi

  echo " "
}
