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
      ${this}.prettyPrintContainerNetworkInformation
      return $?
    fi

    ${this}.printValueOfContainerNetworkInformation
  }

  ${this}.isInputValid() {
    local containerNames=( $(docker container ls --format '{{.Names}}') )
    [[ " ${containerNames[@]} " =~ " ${inputDockerContainer} " ]]
  }

  ${this}.throwInvalidInputError() {
    if [ -z "$inputDockerContainer" ]; then
      echo "${ZSB_ERROR} You have to provide a running docker container"
    else
      echo "${ZSB_ERROR} $(hl ${inputDockerContainer}) is not in the running containers list."
    fi
    return 1
  }

  ${this}.prettyPrintContainerNetworkInformation() {
      ${this}.printContainerNetworkInformation \
        | python3 -m json.tool
  }

  ${this}.printValueOfContainerNetworkInformation() {
    ${this}.printContainerNetworkInformation \
    | ${this}.extractKeyValuePair \
    | ${this}.extractValue \
    | copythis -ps
  }

  ${this}.printContainerNetworkInformation() {
    local NETWORKS_KEY='{{json .NetworkSettings.Networks }}'
    docker inspect ${inputDockerContainer} -f ${NETWORKS_KEY}
  }

  ${this}.extractKeyValuePair() {
    local input=$(< /dev/stdin)
    printf "$input" | grep -Po "${keyRegex[$inputFlag]}" | tail -1 | sed s/,$//
  }

  ${this}.extractValue() {
    local input=$(< /dev/stdin)
    printf "$input" | awk -F '"' '{print $4}'
  }

  ${this}.main "$@"
)

_${zsb}.dkni() {
  case $COMP_CWORD in
    1)
      COMPREPLY=( $(compgen -C "docker ps --format "{{.Names}}"") )
      ;;
    2)
      COMPREPLY=( ip gateway )
      ;;
  esac
}

complete -F _${zsb}.dkni dkni
