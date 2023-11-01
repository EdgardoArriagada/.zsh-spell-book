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

    local -A replaces=(
      ["`wcase -w $file --kebab`"]="`wcase -w $replace --kebab`"
      ["`wcase -w $file --camel`"]="`wcase -w $replace --camel`"
      ["`wcase -w $file --pascal`"]="`wcase -w $replace --pascal`"
    )

    for key value in ${(@kv)replaces}; do
      searchAndReplace $key $value -y
      searchAndReplace -f $key $value -y
    done

  )

  tree $replace
}
