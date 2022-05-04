.github.getNormalUrlFromSshUrl() {
  local input="${1:4}"
  local formattedInput=`printf "$input" | sed 's/:/\//g' | sed 's/\.git$//g'`
  printf "https://${formattedInput}"
}

github() {
  zparseopts -D -E -F -- -compare=compare  || return 1
  local remoteUrl=`git config --get remote.origin.url`

  local url
  if ${zsb}.isUrl "$remoteUrl"; then
    url="$remoteUrl"
  else
    url=`.${0}.getNormalUrlFromSshUrl "$remoteUrl"`
  fi

  local tail
  [[ -n "$compare" ]] && tail="/compare/$(git branch --show-current)"

  zsb_open "${url}${tail}"
}

hisIgnore 'github'

compdef "_${zsb}.nonRepeatedList --compare" github
