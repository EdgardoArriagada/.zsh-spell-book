morning() {
  let attempts=0

  # Unlock sudo password
  sudo cat /dev/null

  let start=`date +%s`
  while true; do
    sudo apt update && sudo apt dist-upgrade -y && break
    : $((attempts++))
    ${zsb}.info "Attempts: ${attempts}"
    sleep 2
  done
  let end=`date +%s`

  a "Morning completed in $((end-start)) seconds."
}

_${zsb}.nocompletion morning

