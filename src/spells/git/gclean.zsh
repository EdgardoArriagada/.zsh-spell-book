gclean() {
  if [[ -n "$1" ]]; then
    git clean -fd "$@"; return $?
  fi

  local untrackedFiles=( $(${zsb}.getGitFiles 'untracked') )

  if [[ -z "$untrackedFiles" ]]; then
    ${zsb}.info "There are no untracked files/directories."; return 0
  fi

  local formattedFiles=$(print -rl -- "  ${(z)^untrackedFiles}")

  ${zsb}.fullPrompt \
    "All untracked files/directories will be deleted." \
    "$formattedFiles"

  ${zsb}.confirmMenu && git clean -fd
}

compdef "_${zsb}.gitUnrepeat 'untracked'" gclean

