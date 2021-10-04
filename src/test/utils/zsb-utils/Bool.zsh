describe "Bool function"; () {
  local -A testData=(
    [0]=true
    ['0']=true
    [" "]=true
    [foo]=true
    ["bar"]=true
    [true]=true

    [1]=false
    ['1']=false
    [""]=false
    [false]=false
  )

  for expected result in "${(@kv)testData}"; do
    it="${result} for '${expected}'"; () {
      [[ "$(Bool $expected)" = $result ]]
      expect $?
    }
  done
}

