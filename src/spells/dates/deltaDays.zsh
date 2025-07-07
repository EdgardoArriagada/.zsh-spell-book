# input in YYYY-MM-DD format
# example: `deltaDays 2022-12-03 2022-12-09`
# or with single date: `deltaDays 2022-12-09` (calculates days from today)

deltaDays() {
  case $# in
    1)
      local date1=`date -j -f '%Y-%m-%d' $1 +%s`
      local date2=`date +%s`
      ;;
    2)
      local date1=`date -j -f '%Y-%m-%d' $1 +%s`
      local date2=`date -j -f '%Y-%m-%d' $2 +%s`
      ;;
    *)
      echo "Usage: deltaDays [date1] date2 (YYYY-MM-DD format)"
      echo "  deltaDays 2022-12-03 2022-12-09  # days between two dates"
      echo "  deltaDays 2022-12-09             # days from today to date"
      return 1
      ;;
  esac

  <<< $(( (date2-date1)/86400 ))
}
