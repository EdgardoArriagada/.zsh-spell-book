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

  alert "DONE: '${packageManager} ${@}'"
}

_${zsb}.run() {
  [[ -f package.json ]] || return 0
  local completion=( `<package.json jq '.scripts' | jq -r 'keys[]'` )

  compadd ${completion}
}


compdef _${zsb}.run run

