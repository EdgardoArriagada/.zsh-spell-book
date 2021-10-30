stopRunningContainers() {
  local runningContainers=$(docker container ls -q)

  [[ -z "$runningContainers" ]] &&
    ${zsb}.info "There are no running containers." &&
    return 0

  ${zsb}.info "Stopping all running containers..."

  docker container stop $(echo "$runningContainers") && \
    ${zsb}.success "All containers did halt."
}

_${zsb}.nocompletion stopRunningContainers

