selectContainer() {
  local inputContainer
  if [[ -z $1 ]]; then
    local runningContainers=$(docker ps --format "{{.Names}}")
    [[ -z $runningContainers ]] && return 1

    inputContainer=$(fzf <<< $runningContainers)
  else
    inputContainer=$1
  fi
  <<< $inputContainer
}

dkattach() {
  local inputContainer=$(selectContainer $1)
  [[ -z "$inputContainer" ]] && ${zsb}.throw "No container"
  printAndRun "docker container attach $inputContainer"
}

dkssh() {
  local inputContainer=$(selectContainer $1)
  [[ -z "$inputContainer" ]] && ${zsb}.throw "No container"
  printAndRun "docker container exec -ti $inputContainer sh"
}

compdef "_${zsb}.singleCompC 'docker ps --format \"{{.Names}}\"'" dkattach dkssh
