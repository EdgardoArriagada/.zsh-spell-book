describe "formatBranch function" $0; () (
  local -A testData=(
    ['feature/AB-123_some-feature']='[ Feature AB-123 ] Some feature'
  )

  for input expected in "${(@kv)testData}"; do
    it="should return '${expected}' for '${input}'"; () {
      local result="$(${zsb}.formatBranch $input)"
      [[ "$result" == "$expected" ]]
      expect $?
    }
  done
)
