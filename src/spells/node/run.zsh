run() {
  local packageManager
  if [[ -f package-lock.json ]]; then
    packageManager="npm run"
  elif [[ -f yarn.lock ]]; then
    packageManager="yarn"
  else
    ${zsb}.cancel 'No package manager found'
  fi

  printAndRun "${packageManager} ${@}"


  if (($#))
    then alert "DONE: '${packageManager} ${@}'"
  fi
}

_${zsb}.run() {
  [[ -f package.json ]] || return 0
  local token='=_=_='
  local parsed=` 
    <package.json jq '.scripts' |\
    jq -r "to_entries|map(\"\(.key)${token}\(.value|tostring)\")|.[]" |\
    sd ':' '\:' |\
    sd "${token}" ':'
  `
  local completion=( ${(@f)parsed} )

  _describe 'command' completion
}


compdef _${zsb}.run run

