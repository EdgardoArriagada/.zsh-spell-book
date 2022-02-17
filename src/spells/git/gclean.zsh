gclean() {
  [[ -n "$1" ]] && git clean -fd "$@" && return $?

  local untrackedFiles=( $(${zsb}.getGitFiles 'untracked') )

  [[ -z "$untrackedFiles" ]] &&
    ${zsb}.cancel "There are no untracked files/directories."

  local formattedFiles=$(print -rl -- "  ${(z)^untrackedFiles}")

  ${zsb}.confirmMenu.withItems \
    "All untracked files/directories will be deleted." \
    "$formattedFiles"

  git clean -fd && ${zsb}.gitStatus
}

compdef "_${zsb}.gitUnrepeat 'untracked'" gclean

