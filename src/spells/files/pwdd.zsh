pwdd() {
  if [[ -n "$1" ]]; then
    print $(print -P %~)/${1}
  else
    print -P %~
  fi
}
