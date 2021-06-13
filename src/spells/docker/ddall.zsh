# like `docker-compose down` but for all the containers
ddall() {
  stopRunningContainers; \
    removeStoppedContainers; \
    docker network prune --force
}

_${zsb}.nocompletion ddall

