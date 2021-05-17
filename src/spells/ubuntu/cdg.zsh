cdg() {
  local DIR="$1"
  if [[ -z "$DIR" ]]; then
    echo "${ZSB_ERROR} You must specify a dir name"
    return 1
  fi

  local FOUND_DIR=$(fastFind $(pwd) "$DIR" d)

  if [[ ! -d "$FOUND_DIR" ]]; then
    echo "${ZSB_INFO} the dir $(hl "$DIR") was not found"
    return 1
  fi

  echo "${ZSB_INFO} Dir found: $(hl "$FOUND_DIR")"
  printf %"$COLUMNS"s | tr " " "_"
  cd "$FOUND_DIR"
}
