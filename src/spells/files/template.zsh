template() {
  local file=${1:?Error: File is required.}

  if [[ ! -f "$file" && ! -d "$file" ]]
    then ${zsb}.throw "File `hl $file` does not exist"
  fi

  tree $file

  printf "\n"

  local replace=$file
  ${zsb}.prompt "Replace with:"
  vared replace

  if [[ -z "$replace" || "$replace" == "$file" ]]
    then ${zsb}.throw "Replacement is required"
  fi

  cp -r $file $replace

  (
    cd $replace

    # wcase -h
    local cases=(
      'flat' 'upper' 'camel' 'pascal' 'snake' 'all-caps' 'kebab' 'train' 'spaced' 'http-header' 'title'
    )

    local -A replaces
    for kase in $cases; do
      replaces[`wcase -w $file --${kase}`]=`wcase -w $replace --${kase}`
    done

    for key value in ${(@kv)replaces}; do
      searchAndReplace $key $value -y
      searchAndReplace -f $key $value -y
    done

  )

  tree $replace
}
