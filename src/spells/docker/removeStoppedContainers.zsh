removeStoppedContainers() {
  local stoppedContainers=$(docker container ls -aq)

  if [ -z "$stoppedContainers" ]; then
    echo "${ZSB_INFO} There are no stopped containers to remove."
    return 0
  fi

  echo "${ZSB_INFO} Removing all stopped Containers..."

  docker container rm $(echo "$stoppedContainers") && \
    echo "${ZSB_SUCCESS} All containers were removed." && \
    return 0
}

complete removeStoppedContainers

