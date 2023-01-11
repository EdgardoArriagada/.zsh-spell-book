# when $# > 1, $1 is a pattern or a comma separated list of blobs
dondedice() {
  case $# in
    0) return 1 ;;
    1) rg -g '!{package-lock.json}' ${1} ;;
    *) rg -g "!{package-lock.json,${1}}" "${@:2}" ;;
  esac
}
