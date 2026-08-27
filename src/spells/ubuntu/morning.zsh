morning() {
  if (( $ZSB_MACOS )); then
    mw
    printAndRun 'brew update && brew upgrade'
    mf
    return $?
  fi

  local attempts=0

  # Unlock sudo password
  sudo cat /dev/null

  mw
  local start=`date +%s`
  while true; do
    sudo apt update && sudo apt dist-upgrade -y && break
    : $((attempts++))
    ${zsb}.info "Attempts: $attempts"
    sleep 2
  done
  local end=`date +%s`

  a "Morning completed in $((end - start)) seconds."
  mf
}

_${zsb}.nocompletion morning


