printColours() (
  integer cols=0
  for i in {0..255}; do
    (( cols++ ))
    local endLine=''
    if (( cols % 10 == 0 )); then endLine="\n"; fi
    printf "\x1b[38;5;${i}m████\x1b[0m►${i}${endLine} "
  done
)

