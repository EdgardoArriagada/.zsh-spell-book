ggl() {
  if [ -z "$1" ]; then
    echo "${ZSB_ERROR} You must provide a branch to pull from."
    return 1
  fi

  git pull origin "$1" && copythis " " -c
}

complete -C "${zsb}_isGitRepo && git branch | sed 's/^\*//'" ggl
