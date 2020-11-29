dkclear() (
  local this="$0"

  ${this}.main() {
    ${this}.printPrompt
    ${this}.playOptionsMenu
  }

  ${this}.printPrompt() {
    echo "${ZSB_WARNING} The following items will be removed"
    echo " "
    echo "$(hl "All running/stopped containers")"
    echo "$(hl "All networks")"
    echo " "
    echo "${ZSB_PROMPT} Do you really want to proceed [Y/n]"
  }

  ${this}.playOptionsMenu() {
    ${zsb}_yesNoMenu ${this}.performClear
  }

  ${this}.performClear() {
    ${this}.stopRunningContainers
    ${this}.pruneContainersVolumesAndNetworks
  }

  ${this}.pruneContainersVolumesAndNetworks() {
    printAndRun 'docker container prune --force' &&
    printAndRun 'docker network prune --force'
  }

  ${this}.main "$@"
)

complete dkclear
