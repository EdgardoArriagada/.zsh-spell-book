removeStoppedContainers() {
  local stoppedContainers=$(docker container ls -aq)

  if [[ -z "$stoppedContainers" ]]; then
    ${zsb}.info "There are no stopped containers to remove."
    return 0
  fi

  ${zsb}.info "Removing all stopped Containers..."

  docker container rm $(echo "$stoppedContainers") && \
    ${zsb}.success "All containers were removed." && \
    return 0
}

complete removeStoppedContainers

