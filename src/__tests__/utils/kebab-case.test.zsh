describe "kebab-case function"; () (
  local -A testData=(
    ["space separated case"]="space-separated-case"
    ['pascalCase']="pascal-case"
    ['kebab-case']="kebab-case"
    ['snake_case']="snake-case"
    ['_bad_snake_case']="bad-snake-case"
    ['.dot-project']=".dot-project"
    ['Capital Letters']="capital-letters"
    ['UpperCamelCase']="upper-camel-case"
    ['multiple  spaces   separated ']="multiple-spaces-separated"
  )

  for input expected in "${(@kv)testData}"; do
    it="should return '${expected}' for '${input}'"; () {
      [[ "$(kebab-case $input)" == "$expected" ]]
      expect $?
    }
  done
)

