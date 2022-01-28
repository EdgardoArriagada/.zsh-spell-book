heimdall() {
  zparseopts -D -E -F -- v=vim || return 1

  local arg="$1"

  if [[ -z "$arg" ]]; then
    print -z "${0} `.heimdall.chooseRequest | fzf`"; return $?
  fi

  if [[ -n "$vim" ]]; then
    eval "nnvim ~/heimdall/${arg}.zsh"; return 0
  fi

  local response=`source ~/heimdall/${arg}.zsh`
  [[ -z "$response" ]] && return 1

  # avoid printing hr output when redirecting to a file
  hr1 >&2

  # if not a valid json
  if ! jq -e . >/dev/null 2>&1 <<<"$response"; then
    print "$response"; return 1
  fi

  jq <<<"$response"
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

