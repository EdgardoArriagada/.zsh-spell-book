${zsb}_getNonRepeatedItems() {
  local inputArray=("$@")

  local uniques=($(echo "${inputArray[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' '))

  declare -A duplicates=( )
  for item in "${inputArray[@]}"; do
    if [ -z "${duplicates[$item]}" ]; then
      duplicates[$item]=0
     else
      (( duplicates[$item] += 1 ))
    fi
  done

  local non_repeated=( )
  for item in "${uniques[@]}"; do
    if [ "${duplicates[$item]}" = "0" ]; then
      non_repeated+=( "$item" )
    fi
  done

  echo ${non_repeated[@]}
}
