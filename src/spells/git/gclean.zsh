gclean() {
  [[ -n "$1" ]] && git clean -fd "$@" && return $?

  local untrackedFiles=( $(${zsb}.getGitFiles 'untracked') )

  [[ -z "$untrackedFiles" ]] &&
    ${zsb}.cancel "There are no untracked files/directories."

  local formattedFiles=$(print -rl -- "  ${(z)^untrackedFiles}")

  ${zsb}.fullPrompt \
    "All untracked files/directories will be deleted." \
    "$formattedFiles"

  ${zsb}.confirmMenu && git clean -fd
}

compdef "_${zsb}.gitUnrepeat 'untracked'" gclean

