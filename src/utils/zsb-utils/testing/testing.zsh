expect() {
  if [[ "$1" = 0 ]]; then
    ${zsb}.pass "$it"
  else
    failingTests+=("$(${zsb}.fail "$it")\n")
  fi
}

describe() {echo "🌑 ${1}"}

runTests() (
  failingTests=()

  local testFiles=( ${ZSB_DIR}/src/__tests__/**/*.zsh )
  ${zsb}.sourceFiles $testFiles

  print "$failingTests"
)
