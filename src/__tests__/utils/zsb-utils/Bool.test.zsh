describe "Bool function"; () (
  local -A testData=(
    [0]=true
    ['0']=true
    [" "]=true
    [foo]=true
    ["bar"]=true
    [true]=true
    [12.34]=true

    ['1']=false
    [""]=false
    [false]=false
  )

  for input expected in "${(@kv)testData}"; do
    it="should return '$expected' for '$input'"; () {
      [[ "$(Bool $input)" == $expected ]]
      expect $?
    }
  done

  # Special case not recognize in array
  it="should return 'false' for '1' as number"; () {
    [[ "$(Bool 1)" == false ]]
    expect $?
  }
)

