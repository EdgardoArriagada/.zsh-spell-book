template() {

  local file=${1:?Error: File is required.}

  if [[ ! -f "$file" && ! -d "$file" ]]
    then ${zsb}.throw "File `hl $file` does not exist"
  fi

  tree $file

  printf "\n"

  local replace
  if [[ -n "$2" ]]
    then replace=$2
    else
      replace=$file
      ${zsb}.prompt "Replace with:"
      vared replace
  fi

  if [[ -z "$replace" || "$replace" == "$file" ]]
    then ${zsb}.throw "Replacement is required"
  fi

  # input validation ends

  function shutdown() {
    # make cursor visible
    tput cnorm
  }

  TRAPINT() { shutdown }
  TRAPQUIT() { shutdown }
  TRAPTERM() { shutdown }

  local spinnerChars=(⣾ ⣽ ⣻ ⢿ ⡿ ⣟ ⣯ ⣷)
  local spinLenMinusOne=$(( ${#spinnerChars[@]} - 1))
  local i=0

  spin() {
    printf "\r${spinnerChars[$i]} "
    (( $i == $spinLenMinusOne )) && i=0 || : $((i++))
  }

  printf "  `color 245 'Replacing...'`"

  cp -r $file $replace


  tput civis # make cursor invisible
  spin

  (
    builtin cd $replace

    # wcase -h
    local cases=(
      'flat' 'upper' 'camel' 'pascal' 'snake' 'all-caps' 'kebab' 'train' 'spaced' 'http-header' 'title'
    )

    local -A replaces
    for kase in $cases; do
      replaces[`wcase -w $file --${kase}`]=`wcase -w $replace --${kase}`
    done

    for key value in ${(@kv)replaces}; do
      spin
      searchAndReplace $key $value -ys
      searchAndReplace $key $value -fys
    done

  )

  printf "\r`color 245 'Done!'`         \n"

  tree $replace
}
