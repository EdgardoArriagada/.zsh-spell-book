gclean() {
  if [[ -n "$1" ]]; then
    git clean -fd "$@"; return $?
  fi

  ${zsb}.warning "All the added files and folders will be deleted."
  echo ""
  ${zsb}.prompt "Are you sure? [Y/n]"

  ${zsb}.confirmMenu && git clean -fd
}

compdef "_${zsb}.gitUnrepeat 'untracked'" gclean

