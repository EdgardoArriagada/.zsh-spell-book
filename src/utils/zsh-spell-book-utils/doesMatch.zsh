${zsb}_doesMatch() {
  local input="$1"
  local regexString="$2"
  return $(echo -E - "$input" | grep -E -q "$regexString")
}

${zsb}_isInteger() {
  local INTEGER_REGEX="^-?[0-9]+$"
  return $(${zsb}_doesMatch "$1" "$INTEGER_REGEX")
}

${zsb}_isFloat() {
  local FLOAT_REGEX="^[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?$"
  return $(${zsb}_doesMatch "$1" "$FLOAT_REGEX")
}

