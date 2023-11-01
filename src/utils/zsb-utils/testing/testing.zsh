expect() { (( $1 == 0 )) && ${zsb}.pass "$it" || ${zsb}.fail "$it" }

describe() {
  local testMessage="$1"
  local fileName="${2##*/}"

  echo "\n🌑 $testMessage $(hl $fileName)"
}

runTests() (
  : ${zsb:='zsb'}
  local testFiles=( $ZSB_DIR/src/__tests__/**/*.test.zsh )
  ${zsb}.sourceFiles $testFiles
)
