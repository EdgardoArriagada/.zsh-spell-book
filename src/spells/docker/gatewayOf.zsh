gatewayOf() (
  local inputDockerContainer="$1"

  main() {
    if ! isInputValid; then
      throwInvalidInputError; return $?
    fi

    printGatewayOfContainer
  }

  isInputValid() [ ! -z "$inputDockerContainer" ]

  throwInvalidInputError() {
    echo "${ZSB_ERROR} You have to provide a running docker container"
    return 1
  }

  printGatewayOfContainer() {
    printContainerNetworkInformation \
    | extractGatewayKeyValuePair \
    | extractGatewayValue
  }

  printContainerNetworkInformation() {
    local NETWORKS_KEY='{{json .NetworkSettings.Networks }}'
    docker inspect ${inputDockerContainer} -f ${NETWORKS_KEY}
  }

  extractGatewayKeyValuePair() {
    local input=$(< /dev/stdin)
    local GATEWAY_KEYVALUE_REGEX='"Gateway":.*?[^\\]",'
    echo "$input" | grep -Po ${GATEWAY_KEYVALUE_REGEX} | sed s/,$//
  }

  extractGatewayValue() {
   local input=$(< /dev/stdin)
   echo "$input" | awk -F '"' '{print $4}'
  }

  main "$@"
)

_${zsb}_gatewayOf() {
  [ "$COMP_CWORD" -gt "1" ] && return 0

  COMPREPLY=( $(compgen -C "docker ps --format "{{.Names}}"") )
}

complete -F _${zsb}_gatewayOf gatewayOf
