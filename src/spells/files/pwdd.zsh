pwdd() {
  local result

  if [[ -n "$1" ]]
    then result="`print -P %~`/$1"
    else result="`print -P %~`"
  fi

  zsb_clipcopy $result
  ${zsb}.info "`hl $result` copied"
}

hisIgnore pwdd

