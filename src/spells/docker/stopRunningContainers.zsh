stopRunningContainers() {
  local runningContainers=`docker container ls -q`

  [[ -z "$runningContainers" ]] &&
    ${zsb}.skip 'There are no running containers.' &&
    return 0

  ${zsb}.info 'Stopping all running containers...'

  docker container stop `printf $runningContainers` &&
    ${zsb}.success 'All containers did halt.'
}

_${zsb}.nocompletion stopRunningContainers

