${zsb}_doesMatch() {
  local word="$1"
  local regexString="$2"
  return $(echo -E - "$word" | grep -E -q "$regexString")
}
