dcos() (
  local this="$0"
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
    echo "${ZSB_WARNING} The following containers will stop:"
    echo " "
    echo "$(hl "$activeServices")"
    echo " "
    echo "${ZSB_PROMPT} Do you really want to stop these containers? [Y/n]"
  }

  ${this}.playOptionsMenu() {
    ${zsb}_confirmMenu && ${this}.stopComposedContainers
  }

  ${this}.stopComposedContainers() printAndRun "docker-compose stop"

  ${this}.main "$@"
)

complete dcos
