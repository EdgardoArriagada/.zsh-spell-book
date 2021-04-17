cds() {
  echo " "
  builtin cd $1 && ls

  if ${zsb}.isGitRepo; then
    ${zsb}.fillWithToken '_'
    ${zsb}.gitStatus
    return 0
  fi

  echo " "
}
