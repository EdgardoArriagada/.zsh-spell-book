alias instala='sudo apt install'
alias updatea='sudo apt update'
alias upgradea='sudo apt dist-upgrade -y'
alias remueve='sudo apt purge'
alias autoremueve='sudo apt clean && sudo apt autoclean && sudo apt autoremove'
alias mata='sudo killall -9'

alias escapeswap='setxkbmap -option caps:swapescape'
alias keysbacktonormal='setxkbmap -option'
alias vouembora='ddall; shutdown -h now'
alias ee="exit"

alias cdesktop='cd ~/Desktop'
alias cdocuments='cd ~/Documents'
alias cdownloads='cd ~/Downloads'
alias links="nvim ~/temp/links.zsh"

alias cmdRank="awk '!(\$1==\"\" || match(\$1, /^(-| |\\\)/)) {print \$1}' ~/.zsh_history | sort | uniq -c | sort -nr | head -20"

hisIgnore ee links cmdRank

whoInPort() { printAndRun "sudo lsof -i \":${1}\"" }
killInPort() { sudo kill -9 $(sudo lsof -t -i ":${1}") }
