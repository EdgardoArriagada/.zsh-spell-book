dcos() (
  local activeServices=$(docker-compose ps --services --filter "status=running")
  main() {
    if [ -z "$activeServices" ]; then
      echo "${ZSB_INFO} No services are running for this docker-compose"
      return 0
    fi

    printPrompt
    playOptionsMenu
  }

  printPrompt() {
    echo "${ZSB_WARNING} The following containers will stop:"
    echo " "
    echo "$(hl "$activeServices")"
    echo " "
    echo "${ZSB_PROMPT} Do you really want to stop these containers? [Y/n]"
  }

  playOptionsMenu() {
    ${zsb}_yesNoMenu stopComposedContainers
  }

  stopComposedContainers() echorun docker-compose stop

  main "$@"
)
