stopRunningContainers() {
  local runningContainers=$(docker container ls -q)

  if [[ -z "$runningContainers" ]]; then
    ${zsb}.info "There are no running containers."
    return 0
  fi

  ${zsb}.info "Stopping all running containers..."

  docker container stop $(echo "$runningContainers") && \
    ${zsb}.success "All containers did halt."
}

_${zsb}.nocompletion stopRunningContainers

