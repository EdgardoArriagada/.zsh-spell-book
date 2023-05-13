run() {
  local runCmd=`[[ -f yarn.lock ]] && <<< 'yarn' || <<< 'npm run'`

  printAndRun "${runCmd} ${@}"
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

