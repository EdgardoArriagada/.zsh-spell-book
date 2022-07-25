pwdd() {
  local result

  zparseopts -D -E -F -- -copy=copy || return 1

  if [[ -n "$1" ]]; then
    result="`print -P %~`/${1}"
  else
    result="`print -P %~`"
  fi

  if [[ -n "$copy" ]]; then
    print "$result" | zsb_clipcopy
    ${zsb}.info "`hl ${result}` copied"
  else
    print "$result"
  fi
}

hisIgnore pwdd

complete -W "--copy" pwdd
