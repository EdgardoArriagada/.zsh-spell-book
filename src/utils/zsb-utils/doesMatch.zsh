# 1=input, 2=regex
${zsb}.doesMatch() [[ "$1" =~ "$2" ]]

${zsb}.isInteger() [[ "$1" = <-> ]]

${zsb}.isFloat() {
  local FLOAT_REGEX="^[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?$"
  ${zsb}.doesMatch "$1" "$FLOAT_REGEX"
}

