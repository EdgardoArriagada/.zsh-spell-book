removeStoppedContainers() {
  local stoppedContainers=`docker container ls -aq`

  [[ -z "$stoppedContainers" ]] &&
    ${zsb}.skip 'There are no stopped containers to remove.' &&
    return 0

  ${zsb}.info 'Removing all stopped Containers...'

  docker container rm `printf $stoppedContainers` && \
    ${zsb}.success 'All containers were removed.'
}

_${zsb}.nocompletion removeStoppedContainers

