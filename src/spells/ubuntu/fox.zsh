fox() (
  local this="$0"

  local inputContent=$(${zsb}.nonPiped "$1")
  local leftRange=$(${zsb}.nonPiped "$2" "$1")
  local rightRange=$(${zsb}.nonPiped "$3" "$2")
  local totalArgs=$(${zsb}.nonPiped "$#" "$(($# + 1))")

  ${this}.main() {
    case "$totalArgs" in
      1)
        ${this}.showInputContent ;;
      2)
        ${this}.showCertainLineOfInput ;;
      3)
        ${this}.showRangeLinesOfInput ;;
      *)
        ${zsb}.throw "Povide a file name with: a line number to copy, two numbers as a range or simply nothing."
    esac
  }

  ${this}.showInputContent() {
    if [[ -f "$inputContent" ]]; then
      cat ${inputContent}
    else
      plainPrint "$inputContent"
    fi
  }

  ${this}.showCertainLineOfInput() {
    if [[ -f "$inputContent" ]]; then
      awk "FNR==$leftRange" ${inputContent}
    else
      plainPrint "$inputContent" | awk "FNR==$leftRange"
    fi
  }

  ${this}.showRangeLinesOfInput() {
    if [[ -f "$inputContent" ]]; then
      awk "FNR>=$leftRange && FNR<=$rightRange" ${inputContent}
    else
      plainPrint "$inputContent" | awk "FNR>=$leftRange && FNR<=$rightRange"
    fi
  }

  ${this}.main "$@"
)

ccp() { fox ${*} | clipcopy }

