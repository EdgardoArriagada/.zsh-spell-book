cds() {
  echo " "
  cd $1

  if ${zsb}.isGitRepo; then
    hr
    ${zsb}.gitStatus
    return 0
  fi

  echo " "
}
