# Docker Network Information
dkni() (
  local this="$0"
  local inputDockerContainer="$1"
  local inputFlag="$2"

  declare -A keyRegex
  keyRegex[ip]='"IPAddress":.*?[^\\]",'
  keyRegex[gateway]='"Gateway":.*?[^\\]",'

  ${this}.main() {
    if ! ${this}.isInputValid; then
      ${this}.throwInvalidInputError; return $?
    fi

    if [ -z "$inputFlag" ]; then
      ${this}.printContainerNetworkInformation
      return $?
    fi

    ${this}.printValueOfContainerNetworkInformation
  }

  ${this}.isInputValid() [ ! -z "$inputDockerContainer" ]

  ${this}.throwInvalidInputError() {
    echo "${ZSB_ERROR} You have to provide a running docker container"
    return 1
  }

  ${this}.printValueOfContainerNetworkInformation() {
    ${this}.printContainerNetworkInformation \
    | ${this}.extractKeyValuePair \
    | ${this}.extractValue
  }

  ${this}.printContainerNetworkInformation() {
    local NETWORKS_KEY='{{json .NetworkSettings.Networks }}'
    docker inspect ${inputDockerContainer} -f ${NETWORKS_KEY}
  }

  ${this}.extractKeyValuePair() {
    local input=$(< /dev/stdin)
    echo "$input" | grep -Po "${keyRegex[$inputFlag]}" | sed s/,$//
  }

  ${this}.extractValue() {
   local input=$(< /dev/stdin)
   echo "$input" | awk -F '"' '{print $4}'
  }

  ${this}.main "$@"
)

_${zsb}.dkni() {
  case $COMP_CWORD in
    1)
      COMPREPLY=( $(compgen -C "docker ps --format "{{.Names}}"") )
      ;;
    2)
      COMPREPLY=( $(compgen -W "ip gateway") )
      ;;
  esac
}

complete -F _${zsb}.dkni dkni
