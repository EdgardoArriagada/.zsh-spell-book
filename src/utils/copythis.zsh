copythis() (
  local this="$0"
  local inputText

  ${this}.main() {
    zparseopts -D -E -F -- s=stayInConsole p=printSuccess || return 1
    inputText="$@"

    ${this}.executeCopyCommand
    ${this}.executeSideEffects
  }

  ${this}.executeCopyCommand() {
    echo -En - ${inputText} | xclip -selection clipboard
  }

  ${this}.executeSideEffects() {
    [[ -n "$printSuccess" ]] && ${zsb}.success "$(hl ${inputText}) copied!"
    [[ -z "$stayInConsole" ]] && close
  }

  ${this}.main "$@"
)

compdef "_${zsb}.nonRepeatedListD '-p:Print if success' '-s:Stay in console'" copythis
