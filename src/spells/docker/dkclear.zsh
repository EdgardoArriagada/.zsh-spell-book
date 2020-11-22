dkclear() (
  main() {
    printPrompt
    playOptionsMenu
  }

  printPrompt() {
    echo "${ZSB_WARNING} The following items will be removed"
    echo " "
    echo "$(hl "All running/stopped containers")"
    echo "$(hl "All networks")"
    echo " "
    echo "${ZSB_PROMPT} Do you really want to proceed [Y/n]"
  }

  playOptionsMenu() {
    ${zsb}_yesNoMenu performClear
  }

  performClear() {
    stopRunningContainers
    pruneContainersVolumesAndNetworks
  }

  pruneContainersVolumesAndNetworks() {
    printAndRun 'docker container prune --force' &&
    printAndRun 'docker network prune --force'
  }

  main "$@"
)

complete dkclear
