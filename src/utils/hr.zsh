hr() {
  local token=${1:='·'}
  printf %"$COLUMNS"s | tr " " "$token"
}

hr1() color 4 `hr`
