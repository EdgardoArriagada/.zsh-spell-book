expect() {
 [[ "$1" = 0 ]] && ${zsb}.pass "$it" || ${zsb}.fail "$it"
}

describe() {echo "🌑 ${1}"}
