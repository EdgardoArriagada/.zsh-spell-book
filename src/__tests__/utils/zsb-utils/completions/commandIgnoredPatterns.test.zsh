describe "command completion ignored patterns" $0; () (
  it="should ignore underscore commands when the prefix does not start with underscore"; () {
    local PREFIX="vno"
    local -a reply

    _${zsb}.commandIgnoredPatterns
    local commandStatus=$?

    [[ $commandStatus -eq 0 && "${(j: :)reply}" == "_*" ]]
    expect $?
  }

  it="should allow underscore commands when the prefix starts with underscore"; () {
    local PREFIX="_v"
    local -a reply

    _${zsb}.commandIgnoredPatterns
    local commandStatus=$?

    [[ $commandStatus -eq 0 && ${#reply[@]} -eq 0 ]]
    expect $?
  }
)
