template() {
  local file=${1:?File is required}

  if [[ ! -f "$file" && ! -d "$file" ]]
    then ${zsb}.throw "File `hl ${file}` does not exist"
  fi

  tree $file

  printf "\n"

  local replace=$file
  ${zsb}.prompt "Replace with:"
  vared replace

  if [[ -z "$replace" || "$replace" == "$file" ]]
    then ${zsb}.throw "Replacement is required"
  fi

  cp -r ${file} ${replace}

  (
    cd $replace

    local replace_u=`firstToUpper ${replace}`
    local replace_l=`firstToLower ${replace}`
    local file_u=`firstToUpper ${file}`
    local file_l=`firstToLower ${file}`

    searchAndReplace ${file_u} ${replace_u} -y
    searchAndReplace ${file_l} ${replace_l} -y

    searchAndReplace -f ${file_u} ${replace_u} -y
    searchAndReplace -f ${file_l} ${replace_l} -y
  )

  tree $replace
}
