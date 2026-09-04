describe "co command" $0; () (
  codex() { calls+=( "${(j: :)@}" ) }

  local -a calls

  it="should dispatch supported invocations"; () {
    co
    co --update-pr-title-and-description
    co --code-review
    co --commit

    [[ "${(j:|:)calls}" == '--dangerously-bypass-approvals-and-sandbox|exec Update this pr title and description|$thermo-nuclear-code-quality-review|exec Create a Git commit from the currently staged files' ]]
    expect $?
  }

  it="should reject unsupported invocations"; () {
    local unknownStatus multipleStatus
    co --unknown 2>/dev/null
    unknownStatus=$?
    co --code-review extra 2>/dev/null
    multipleStatus=$?

    [[ $unknownStatus -ne 0 && $multipleStatus -ne 0 ]]
    expect $?
  }

  it="should complete supported options"; () {
    _describe() { completions=( "${(@P)2}" ) }

    local CURRENT=2
    local -a completions
    _${zsb}.co

    [[ "${(j:|:)completions}" == '--update-pr-title-and-description:update PR title and description|--code-review:run thermo-nuclear code quality review|--commit:create a git commit' ]]
    expect $?
  }
)
