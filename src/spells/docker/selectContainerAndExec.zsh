${zsb}.selectContainerAndExecute() {
  local rawCmd=$1
  shift

  local inputContainer
  if [[ -z $1 ]]; then
    local runningContainers=$(docker ps --format "{{.Names}}")
    [[ -z $runningContainers ]] && ${zsb}.throw "No containers found"

    inputContainer=$(fzf <<< $runningContainers)
  else
    inputContainer=$1
  fi

  local cmd=$(echo $rawCmd | sed "s/{inputContainer}/$inputContainer/g")
  printAndRun $cmd
}

dkattach() ${zsb}.selectContainerAndExecute 'docker container attach {inputContainer}' $1
dkssh() ${zsb}.selectContainerAndExecute 'docker container exec -ti {inputContainer} sh' $1

compdef "_${zsb}.singleCompC 'docker ps --format \"{{.Names}}\"'" dkattach dkssh
