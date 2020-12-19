${zsb}.dockerComposeHalt() (
  local this="$0"
  local callback="$1"

  ${this}.main() {
    ${this}.printPrompt
    ${this}.playOptionsMenu
  }

  ${this}.printPrompt() {
    echo "${ZSB_PROMPT} Do you really want to excecute $(hl ${callback})? [Y/n]"
  }

  ${this}.playOptionsMenu() {
    ${zsb}.confirmMenu && ${this}.haltComposedContainers
  }

  ${this}.haltComposedContainers() printAndRun "${callback}"

  ${this}.main "$@"
)

alias dcos="${zsb}.dockerComposeHalt 'docker-compose stop'"
alias dcod="${zsb}.dockerComposeHalt 'docker-compose down'"

complete ${zsb}.dockerComposeHalt

