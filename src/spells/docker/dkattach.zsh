dkattach() {
  local inputContainer
  if [[ -z $1 ]]; then
    local runningContainers=$(docker ps --format "{{.Names}}")
    [[ -z $runningContainers ]] && ${zsb}.throw "No running containers found."

    inputContainer=$(fzf <<< $runningContainers)
  else
    inputContainer=$1
  fi

  printAndRun "docker container attach $inputContainer"
}

compdef "_${zsb}.singleCompC 'docker ps --format \"{{.Names}}\"'" dkattach

