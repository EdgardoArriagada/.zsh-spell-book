gatewayOf() (
  local this="$0"
  local inputDockerContainer="$1"

  ${this}.main() {
    if ! ${this}.isInputValid; then
      ${this}.throwInvalidInputError; return $?
    fi

    ${this}.printGatewayOfContainer
  }

  ${this}.isInputValid() [ ! -z "$inputDockerContainer" ]

  ${this}.throwInvalidInputError() {
    echo "${ZSB_ERROR} You have to provide a running docker container"
    return 1
  }

  ${this}.printGatewayOfContainer() {
    ${this}.printContainerNetworkInformation \
    | ${this}.extractGatewayKeyValuePair \
    | ${this}.extractGatewayValue
  }

  ${this}.printContainerNetworkInformation() {
    local NETWORKS_KEY='{{json .NetworkSettings.Networks }}'
    docker inspect ${inputDockerContainer} -f ${NETWORKS_KEY}
  }

  ${this}.extractGatewayKeyValuePair() {
    local input=$(< /dev/stdin)
    local GATEWAY_KEYVALUE_REGEX='"Gateway":.*?[^\\]",'
    echo "$input" | grep -Po ${GATEWAY_KEYVALUE_REGEX} | sed s/,$//
  }

  ${this}.extractGatewayValue() {
   local input=$(< /dev/stdin)
   echo "$input" | awk -F '"' '{print $4}'
  }

  ${this}.main "$@"
)

_${zsb}.gatewayOf() {
  [ "$COMP_CWORD" -gt "1" ] && return 0

  COMPREPLY=( $(compgen -C "docker ps --format "{{.Names}}"") )
}

complete -F _${zsb}.gatewayOf gatewayOf
