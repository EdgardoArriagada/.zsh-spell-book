cds() {
  echo " "
  builtin cd $1 && ls

  if isGitRepo; then
    printf %"$COLUMNS"s | tr " " "_"
    gitStatus
    return 0
  fi

  echo " "
}
