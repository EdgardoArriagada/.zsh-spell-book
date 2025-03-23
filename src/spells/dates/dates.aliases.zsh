${zsb}.cpdate_generic() { zsb_clipcopy $1 && ${zsb}.success "`hl $1` as `hl $2` copied to clipboard."; }

cpdate() {
  zparseopts -D -E -F -- \
    f1=format1 \
    f2=format2 \
    || return 1

  if [[ -n "$format1" ]]; then
    ${zsb}.cpdate_generic `date +%Y-%m-%d` YYYY-MM-DD
  elif [[ -n "$format2" ]]; then
    ${zsb}.cpdate_generic `date +%d-%m-%Y` DD-MM-YYYY
  else
    ${zsb}.cpdate_generic `date +%Y-%m-%d` YYYY-MM-DD
  fi
}

compdef "_${zsb}.singleComp \
  '-f1:YYYY-MM-DD' \
  '-f2:DD-MM-YYYY' \
  " cpdate
