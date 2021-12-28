cds() {
  echo " "
  cd $1

  if ${zsb}.isGitRepo; then
    ${zsb}.fillWithToken '_'
    ${zsb}.gitStatus
    return 0
  fi

  echo " "
}
