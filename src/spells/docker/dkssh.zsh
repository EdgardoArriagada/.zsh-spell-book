dkssh() {
  local inputContainer=${1?Error: You must provide a container.}
  printAndRun "docker container exec -ti ${inputContainer} bash"
}

compdef "_${zsb}.singleCompC 'docker ps --format \"{{.Names}}\"'" dkssh

