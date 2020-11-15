ggl() {
  if [ -z "$1" ]; then
    echo "${ZSB_ERROR} You must provide a branch to pull from."
    return 1
  fi

  git pull origin "$1" && copythis " " -c
}

complete -C "$ZSB_GIT_BRANCHES" ggl
