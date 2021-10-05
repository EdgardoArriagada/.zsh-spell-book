expect() { (( $1 == 0 )) && ${zsb}.pass "$it" || ${zsb}.fail "$it" }

describe() {echo "\n🌑 ${1}"}

runTests() (
  local testFiles=( ${ZSB_DIR}/src/__tests__/**/*.test.zsh )
  ${zsb}.sourceFiles $testFiles
)
