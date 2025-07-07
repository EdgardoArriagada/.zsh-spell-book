# input in YYYY-MM-DD format
# example: `deltaDays 2022-12-03 2022-12-09`
# or with single date: `deltaDays 2022-12-09` (calculates days from today)

deltaDays() {
  if [[ $# -eq 1 ]]; then
    # Single date provided - use today as date1, provided date as date2
    local date2=`date +%s`
    local date1=`date -j -f '%Y-%m-%d' $1 +%s`
  elif [[ $# -eq 2 ]]; then
    # Two dates provided - use original behavior
    local date1=`date -j -f '%Y-%m-%d' $1 +%s`
    local date2=`date -j -f '%Y-%m-%d' $2 +%s`
  else
    echo "Usage: deltaDays [date1] date2 (YYYY-MM-DD format)"
    echo "  deltaDays 2022-12-03 2022-12-09  # days between two dates"
    echo "  deltaDays 2022-12-09             # days from today to date"
    return 1
  fi

  <<< $(( (date2-date1)/86400 ))
}
