github() {
  zparseopts -D -E -F -- -compare=compare  || return 1

  ${zsb}.assertGitRepo

  local url=`get_repo_url`

  if [[ -z "$url" ]]
    then ${zsb}.throw "Remote url not set."
  fi

  if [[ -n "$compare" ]]
    then zsb_open "$url/compare/`git branch --show-current`"
    else zsb_open $url
  fi
}

hisIgnore 'github'

compdef "_${zsb}.nonRepeatedList --compare" github
