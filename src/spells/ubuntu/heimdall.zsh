heimdall() {
  zparseopts -D -E -F -- v=vim || return 1

  local arg="$1"

  if [[ -z "$arg" ]]; then
    print -z "${0} `.heimdall.chooseRequest | fzf`"
    return $?
  fi

  if [[ -n "$vim" ]]; then
    eval "nnvim ~/heimdall/${arg}.zsh"
    return 0
  fi

  local response=`source ~/heimdall/${arg}.zsh`
  hr >&2

  if jq -e . >/dev/null 2>&1 <<<"$response"; then
    jq <<<"$response"
  else
    print "$response"
  fi
}

.heimdall.chooseRequest() {
 fd '.zsh$' ~/heimdall \
 | sed -e 's/^.*heimdall\///' \
 | sed -e 's/\.zsh$//'
}

alias cdheimdall='cds ~/heimdall'

_${zsb}.heimdall() {
  local newCompletion=( `.heimdall.chooseRequest` )
  _describe 'command' newCompletion
}

compdef _${zsb}.heimdall heimdall

