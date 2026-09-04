${zsb}.cpdate_generic() { zsb_clipcopy "$1" && ${zsb}.success "`hl "$1"` as `hl "$2"` copied to clipboard."; }

cpdate() {
  zparseopts -D -E -F -- \
    f1=format1 \
    f2=format2 \
    f3=format3 \
    || return 1

  if [[ -n "$format1" ]]; then
    ${zsb}.cpdate_generic "$(command date +%Y-%m-%d)" YYYY-MM-DD
  elif [[ -n "$format2" ]]; then
    ${zsb}.cpdate_generic "$(command date +%d-%m-%Y)" DD-MM-YYYY
  elif [[ -n "$format3" ]]; then
    local -a weekdays=(lunes martes miércoles jueves viernes sábado domingo)
    local -a months=(enero febrero marzo abril mayo junio julio agosto septiembre octubre noviembre diciembre)
    local weekday day month year
    read weekday day month year <<< "$(command date '+%u %d %m %Y')"
    ${zsb}.cpdate_generic "${weekdays[$((10#$weekday))]} $((10#$day)) de ${months[$((10#$month))]} de $year" 'viernes 15 de mayo de 2026'
  else
    ${zsb}.cpdate_generic "$(command date +%Y-%m-%d)" YYYY-MM-DD
  fi
}

compdef "_${zsb}.singleComp \
  '-f1:YYYY-MM-DD' \
  '-f2:DD-MM-YYYY' \
  '-f3:viernes 15 de mayo de 2026' \
  " cpdate
