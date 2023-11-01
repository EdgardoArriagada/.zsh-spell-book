# when $# > 1, $1 is a comma separated list of globs
dondedice() {
  local maxColumns=`${zsb}.getMaxSearchColumns`
  case $# in
    0) return 1 ;;
    1) rg --glob '!{package-lock.json}' $1 -M $maxColumns ;;
    *) rg --glob "!{package-lock.json,${1}}" "${@:2}" -M $maxColumns ;;
  esac
}
