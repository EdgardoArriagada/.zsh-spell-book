# when $# > 1, $1 is a comma separated list of globs
dondedice() {
  case $# in
    0) return 1 ;;
    1) rg --glob '!{package-lock.json}' ${1} ;;
    *) rg --glob "!{package-lock.json,${1}}" "${@:2}" ;;
  esac
}
