hr() {
  local token=${1:='_'}
  printf %"$COLUMNS"s | tr " " "$token"
}
