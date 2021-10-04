describe "Bool function"; () {
  local -A testData=(
    [0]=true
    ['0']=true
    [" "]=true
    [foo]=true
    ["bar"]=true
    [true]=true

    ['1']=false
    [""]=false
    [false]=false
  )

  for expected result in "${(@kv)testData}"; do
    it="should return '${result}' for '${expected}'"; () {
      [[ "$(Bool $expected)" = $result ]]
      expect $?
    }
  done

  # Special case not recognize in array
  it="should return 'false' for '1' as number"; () {
    [[ "$(Bool 1)" = false ]]
    expect $?
  }
}

