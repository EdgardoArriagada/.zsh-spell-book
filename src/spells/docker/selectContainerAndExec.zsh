${zsb}.selectContainerAndExecute() {
  local invoker=$1
  local rawCmd=$2
  shift 2

  local inputContainer=$1
  if [[ -z "$inputContainer" ]]; then
    local runningContainers=$(docker ps --format "{{.Names}}")
    [[ -z "$runningContainers" ]] && ${zsb}.throw "No containers found"

    inputContainer=$(fzf <<< $runningContainers)
    print -z "$invoker $inputContainer"
    return 0
  fi

  local cmd=$(echo $rawCmd | sed "s/{inputContainer}/$inputContainer/g")
  printAndRun $cmd
}

dkattach() ${zsb}.selectContainerAndExecute $0 'docker container attach {inputContainer}' $1
dkssh() ${zsb}.selectContainerAndExecute $0 'docker container exec -ti {inputContainer} sh' $1

compdef "_${zsb}.singleCompC 'docker ps --format \"{{.Names}}\"'" dkattach dkssh
