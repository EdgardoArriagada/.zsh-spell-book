alias instala='sudo apt install'
alias updatea='sudo apt update'
alias upgradea='sudo apt dist-upgrade -y'
alias remueve='sudo apt purge'
alias autoremueve='sudo apt clean && sudo apt autoclean && sudo apt autoremove'
alias mata='sudo killall -9'

alias escapeswap='setxkbmap -option caps:swapescape'
alias keysbacktonormal='setxkbmap -option'
alias vouembora='ddall && shutdown -h now'
alias ee="exit"


morning() {
  local attempts=0

  # Unlock sudo password
  sudo cat /dev/null

  local -r start=`date +%s`
  while true; do
    sudo apt update && sudo apt dist-upgrade -y && break
    : $((attempts++))
    ${zsb}.info "Attempts: ${attempts}"
    sleep 2
  done
  local -r end=`date +%s`

  a "Morning completed in $((end - start)) seconds."
}
