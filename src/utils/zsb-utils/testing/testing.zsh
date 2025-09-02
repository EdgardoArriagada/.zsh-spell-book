expect() { (( $1 == 0 )) && ${zsb}.pass "$it" || ${zsb}.fail "$it" }

describe() {
  local testMessage="$1"
  local filename="${2##*/}"

  echo "\n🌑 $testMessage $(hl $filename)"
}

runTests() (
  : ${zsb:='zsb'}
  local testFiles=( $ZSB_DIR/src/__tests__/**/*.test.zsh )
  for f in $testFiles; do source $f; done
)
