printColours() (
  local this="$0"
  integer cols=0

  ${this}.printRange() {
    local left="$1"
    local right="$2"
    local spaces="$3"

    for i in {${left}..${right}}; do
      (( cols++ ))
      local endLine=''
      if (( cols % 10 == 0 )); then endLine="\n"; fi
      printf " \e[38;5;${i}m████\033[0m►${i}${spaces}${endLine}"
    done
  }

  ${this}.printRange 0 9 '  '
  ${this}.printRange 10 99 ' '
  ${this}.printRange 100 255 ''
)

