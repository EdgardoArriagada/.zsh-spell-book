cds() {
  print " "
  cd $1

  if ${zsb}.isGitRepo; then
    hr1
    ${zsb}.gitStatus
    return 0
  fi

  print " "
}
