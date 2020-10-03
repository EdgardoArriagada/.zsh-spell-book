${zsb}_isInteger() {
  local INTEGER_REGEX="^-?[0-9]+$"
  return $(${zsb}_doesMatch "$1" "$INTEGER_REGEX")
}
