describe "cpdate command" $0; () (
  ${zsb}.cpdate_generic() { calls+=( "$2" ); }
  ${zsb}.info() { :; }

  local -a calls

  it="should dispatch supported formats"; () {
    cpdate
    cpdate -f1
    cpdate -f2
    cpdate -f3

    [[ "${(j:|:)calls}" == 'YYYY-MM-DD|YYYY-MM-DD|DD-MM-YYYY|viernes 15 de mayo de 2026' ]]
    expect $?
  }

  it="should reject unsupported invocations"; () {
    local unknownStatus multipleStatus
    cpdate --unknown 2>/dev/null
    unknownStatus=$?
    cpdate -f1 -f2 2>/dev/null
    multipleStatus=$?

    [[ $unknownStatus -ne 0 && $multipleStatus -ne 0 ]]
    expect $?
  }

  it="should complete supported options"; () {
    _describe() { completions=( "${(@P)2}" ); }

    local CURRENT=2
    local -a completions
    _${zsb}.cpdate

    [[ "${(j:|:)completions}" == '-f1:YYYY-MM-DD|-f2:DD-MM-YYYY|-f3:viernes 15 de mayo de 2026' ]]
    expect $?
  }
)
