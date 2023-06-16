github() {
  zparseopts -D -E -F -- -compare=compare  || return 1

  ${zsb}.validateGitRepo

  local url=`get_repo_url`

  if [[ -z "$url" ]]; then
    ${zsb}.throw "Remote url not set."
  fi

  local tail
  [[ -n "$compare" ]] && tail="/compare/`git branch --show-current`"

  zsb_open "${url}${tail}"
}

hisIgnore 'github'

compdef "_${zsb}.nonRepeatedList --compare" github
