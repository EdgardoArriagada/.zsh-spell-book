doesMatch() {
  local word="$1"
  local regexString="$2"
  return $(echo "$word" | grep -E -q "$regexString")
}
