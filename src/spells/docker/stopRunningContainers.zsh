stopRunningContainers() {
  local runningContainers=$(docker container ls -q)

  if [ -z "$runningContainers" ]; then
    echo "${ZSB_INFO} There are no running containers."
    return 0
  fi

  echo "${ZSB_INFO} Stopping all running containers..."

  docker container stop $(echo "$runningContainers") && \
    echo "${ZSB_SUCCESS} All containers did halt." && \
    return 0
}

complete stopRunningContainers

