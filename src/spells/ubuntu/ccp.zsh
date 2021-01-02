ccp() {
  local inputFile="$1"
  local leftRange="$2"
  local rightRange="$3"

  case "$#" in
    1)
      cat "$inputFile" | clipcopy ;;
    2)
      awk "FNR==$leftRange" $inputFile | clipcopy ;;
    3)
      awk "FNR>=$leftRange && FNR<=$rightRange" $inputFile | clipcopy ;;
    *)
      echo "${ZSB_ERROR} You must provide a file, an optional line to copy or two numbers describing a range of lines to copy."
  esac
}
