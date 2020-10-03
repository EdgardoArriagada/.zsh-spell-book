isInteger() {
  local INTEGER_REGEX="^-?[0-9]+$"
  return $(doesMatch "$1" "$INTEGER_REGEX")
}
