dkclear() (
  main() {
    printPrompt
    playOptionsMenu
  }

  printPrompt() {
    echo "${ZSB_WARNING} The following items will be removed"
    echo " "
    echo "$(hl "All running/stopped containers")"
    echo "$(hl "All volumes")"
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

  stopRunningContainers() {
    local containers=$(docker container ls -aq)
    if [ ! -z "$containers" ]; then
      echo "${ZSB_INFO} Stopping all running containers"
      docker stop $(echo "$containers") && echo "${ZSB_SUCCESS} All running containers stopped\n" &&
    fi
  }

  pruneContainersVolumesAndNetworks() {
    echorun 'docker container prune --force' &&
    echorun 'docker volume prune --force' &&
    echorun 'docker network prune --force'
  }

  main "$@"
)
