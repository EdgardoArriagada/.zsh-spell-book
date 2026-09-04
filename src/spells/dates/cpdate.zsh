${zsb}.cpdate_generic() { zsb_clipcopy "$1" && ${zsb}.success "`hl "$1"` as `hl "$2"` copied to clipboard."; }

cpdate() {
  case "$#:${1-}" in
    0:|1:-f1)
      ${zsb}.cpdate_generic "$(command date +%Y-%m-%d)" YYYY-MM-DD
      ;;
    1:-f2)
      ${zsb}.cpdate_generic "$(command date +%d-%m-%Y)" DD-MM-YYYY
      ;;
    1:-f3)
      local -a weekdays=(lunes martes miércoles jueves viernes sábado domingo)
      local -a months=(enero febrero marzo abril mayo junio julio agosto septiembre octubre noviembre diciembre)
      local weekday day month year
      read weekday day month year <<< "$(command date '+%u %d %m %Y')"
      ${zsb}.cpdate_generic "${weekdays[$((10#$weekday))]} $((10#$day)) de ${months[$((10#$month))]} de $year" 'viernes 15 de mayo de 2026'
      ;;
    *)
      ${zsb}.info 'Usage: cpdate [-f1 | -f2 | -f3]'
      return 1
      ;;
  esac
}

_${zsb}.cpdate() {
  (( CURRENT > 2 )) && return 0

  local -a options=(
    '-f1:YYYY-MM-DD'
    '-f2:DD-MM-YYYY'
    '-f3:viernes 31 de octubre de 2025'
  )
  _describe 'option' options
}

compdef _${zsb}.cpdate cpdate
