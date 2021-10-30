removeStoppedContainers() {
  local stoppedContainers=$(docker container ls -aq)

  [[ -z "$stoppedContainers" ]] &&
    ${zsb}.cancel "There are no stopped containers to remove."

  ${zsb}.info "Removing all stopped Containers..."

  docker container rm $(echo "$stoppedContainers") && \
    ${zsb}.success "All containers were removed."
}

_${zsb}.nocompletion removeStoppedContainers

