gb() {
  if [ -z "$1" ]; then
    echo "\n" && git branch && echo "\n"
    return 0
  fi

  git branch "$@"
  return 0
}

complete -C "$ZSB_GIT_BRANCHES" gb
