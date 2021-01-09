fox() {
  local inputFile="$1"
  local leftRange="$2"
  local rightRange="$3"

  case "$#" in
    1)
      < "$inputFile" ;;
    2)
      awk "FNR==$leftRange" $inputFile ;;
    3)
      awk "FNR>=$leftRange && FNR<=$rightRange" $inputFile ;;
    *)
      echo "${ZSB_ERROR} Povide a file name with: a line number to copy, two numbers as a range or simply nothing."
  esac
}

ccp() {
  fox ${*} | clipcopy
}

