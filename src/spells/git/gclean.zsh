gclean() {
  if [[ -n "$1" ]]; then
    git clean -fd "$@"; return $?
  fi

  local untrackedFiles=( $(${zsb}.getGitFiles 'untracked') )

  if [[ -z "$untrackedFiles" ]]; then
    ${zsb}.info "There are no untracked files."; return 0
  fi

  local formattedFiles=$(print -rl -- "  ${(z)^untrackedFiles}")

  ${zsb}.warning "All untracked files will be deleted."
  echo " "
  echo "$(hl $formattedFiles)"
  echo " "
  ${zsb}.prompt "Are you sure? [Y/n]"

  ${zsb}.confirmMenu && git clean -fd
}

compdef "_${zsb}.gitUnrepeat 'untracked'" gclean

