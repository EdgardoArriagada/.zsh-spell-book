# input in YYYY-MM-DD format
# example: `deltaDays 2022-12-03 2022-12-09`

deltaDays() {
  local date1=`date -j -f '%Y-%m-%d'  ${1} +%s`
  local date2=`date -j -f '%Y-%m-%d'  ${2} +%s`

  printf $(( (date2-date1)/86400 ))
}
