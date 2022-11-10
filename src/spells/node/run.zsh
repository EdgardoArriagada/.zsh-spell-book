run() {
  local runCmd=`[[ -f yarn.lock ]] && print "yarn" || print "npm run"`

  printAndRun "${runCmd} ${@}"

  if (($#))
    then alert "DONE: '${runCmd} ${@}'"
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

