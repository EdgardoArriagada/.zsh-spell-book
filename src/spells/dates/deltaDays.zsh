deltaDays() {
  # input in YYYY-MM-DD format
  # count 2022-12-03 2022-12-09
  date1=$(date -j -f "%Y-%m-%d"  "$1" +%s)
  date2=$(date -j -f "%Y-%m-%d"  "$2" +%s)

  printf '%d delta days\n' "$(( (date2-date1)/86400 ))"
}
