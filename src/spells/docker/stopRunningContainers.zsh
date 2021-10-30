stopRunningContainers() {
  local runningContainers=$(docker container ls -q)

  [[ -z "$runningContainers" ]] &&
    ${zsb}.cancel "There are no running containers."

  ${zsb}.info "Stopping all running containers..."

  docker container stop $(echo "$runningContainers") && \
    ${zsb}.success "All containers did halt."
}

_${zsb}.nocompletion stopRunningContainers

