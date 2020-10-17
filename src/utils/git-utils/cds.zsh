cds() {
  echo " "
  builtin cd $1 && ls

  if ${zsb}_isGitRepo; then
    ${zsb}_printHr '--no-space'
    ${zsb}_gitStatus
    return 0
  fi

  echo " "
}
