describe "command completion styles" $0; () (
  ${zsb}.test.sourceCompletionStyleBlock() {
    zstyle -d ':completion:*' matcher-list
    zstyle -d ':completion:*:*:-command-:*:*' ignored-patterns
    zstyle -d ':completion:*:ignored:-command-:*:*' matcher-list

    local zshConfig="${${(%):-%x}:A:h}/../../../../zsh.config.zsh"

    source <(
      sed -n '/# Make completion smarter/,/^$/p' "$zshConfig"
    )
  }

  it="should ignore underscore commands when the prefix does not start with underscore"; () {
    ${zsb}.test.sourceCompletionStyleBlock

    local PREFIX="vno"
    local -a ignoredPatterns

    zstyle -a ':completion:*:*:-command-:*:*' ignored-patterns ignoredPatterns

    [[ "${(j: :)ignoredPatterns}" == "_*" ]]
    expect $?
  }

  it="should allow underscore commands when the prefix starts with underscore"; () {
    ${zsb}.test.sourceCompletionStyleBlock

    local PREFIX="_v"
    local -a ignoredPatterns

    zstyle -a ':completion:*:*:-command-:*:*' ignored-patterns ignoredPatterns

    [[ ${#ignoredPatterns[@]} -eq 0 ]]
    expect $?
  }

  it="should keep ignored command matches case insensitive"; () {
    ${zsb}.test.sourceCompletionStyleBlock

    local -a matcherList

    zstyle -a ':completion:*:ignored:-command-:*:*' matcher-list matcherList

    [[ "${(j: :)matcherList}" == "m:{[:lower:]}={[:upper:]}" ]]
    expect $?
  }
)
