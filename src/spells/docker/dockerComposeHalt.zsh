${zsb}.dockerComposeHalt() (
  local this="$0"
  local callback="$1"
  local activeServices=$(docker-compose ps --services --filter "status=running")

  ${this}.main() {
    if [ -z "$activeServices" ]; then
      echo "${ZSB_INFO} No services are running for this docker-compose"
      return 0
    fi

    ${this}.printPrompt
    ${this}.playOptionsMenu
  }

  ${this}.printPrompt() {
    echo "${ZSB_WARNING} The following containers will halt with $(hl ${callback}):"
    echo " "
    echo "$(hl "$activeServices")"
    echo " "
    echo "${ZSB_PROMPT} Do you really want to stop these containers? [Y/n]"
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

