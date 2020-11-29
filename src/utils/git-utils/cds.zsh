cds() {
  echo " "
  builtin cd $1 && ls

  if ${zsb}.isGitRepo; then
    ${zsb}.printHr '--no-space'
    ${zsb}.gitStatus
    return 0
  fi

  echo " "
}
