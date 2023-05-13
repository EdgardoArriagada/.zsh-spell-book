ZSB_MESSAGE_BADGE=`color 245 zsb`

${zsb}.message() <<< "${ZSB_MESSAGE_BADGE} `color ${1} ${2}`: ${@:3}"

${zsb}.throw() {
  ${zsb}.message $ZSB_RED 'ERROR' "$@"
  throw Error
}

${zsb}.cancel() {
  ${zsb}.message $ZSB_CYAN 'CANCEL' "$@"
  throw Error
}

${zsb}.success() ${zsb}.message $ZSB_GREEN 'SUCCESS' "$@"

${zsb}.fail() ${zsb}.message $ZSB_RED 'FAIL' "$@"

${zsb}.pass() ${zsb}.message $ZSB_GREEN 'PASS' "$@"

${zsb}.warning() ${zsb}.message $ZSB_YELLOW 'WARNING' "$@"

${zsb}.info() ${zsb}.message $ZSB_BLUE 'INFO' "$@"

${zsb}.prompt() ${zsb}.message $ZSB_PURPLE 'PROMPT' "$@"

${zsb}.expected() ${zsb}.message $ZSB_GREEN '+ Expected' "$@"

${zsb}.current() ${zsb}.message $ZSB_RED '- Current' "$@"

${zsb}.validate() {
  [[ -z "${(P)1}" ]] && ${zsb}.throw "You must set `hl ${1}` first."
  return 0
}

