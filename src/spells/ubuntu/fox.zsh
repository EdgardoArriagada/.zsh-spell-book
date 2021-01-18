fox() {
  local inputContent=$(${zsb}.nonPiped "$1")

  # if is a valid filename
  if [ -f "$inputContent" ]; then
    inputContent=$(< "$inputContent")
  fi

  local leftRange=$(${zsb}.nonPiped "$2" "$1")
  local rightRange=$(${zsb}.nonPiped "$3" "$2")
  local totalArgs=$(${zsb}.nonPiped "$#" "$(($# + 1))")

  case "$totalArgs" in
    1)
      echo "$inputContent" ;;
    2)
      echo "$inputContent" | awk "FNR==$leftRange" ;;
    3)
      echo "$inputContent" | awk "FNR>=$leftRange && FNR<=$rightRange" ;;
    *)
      echo "${ZSB_ERROR} Povide a file name with: a line number to copy, two numbers as a range or simply nothing."
  esac
}

ccp() {
  fox ${*} | clipcopy
}

