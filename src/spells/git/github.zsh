github() {
  zparseopts -D -E -F -- -compare=compare  || return 1

  local url=`get_repo_url`

  local tail
  [[ -n "$compare" ]] && tail="/compare/`git branch --show-current`"

  zsb_open "${url}${tail}"
}

hisIgnore 'github'

compdef "_${zsb}.nonRepeatedList --compare" github
