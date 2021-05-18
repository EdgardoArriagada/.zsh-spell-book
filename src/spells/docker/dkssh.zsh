dkssh() {
  local -r inputContainer=${1:?'You must provide a container'}
  printAndRun "docker container exec -ti ${inputContainer} bash"
}

compdef "_${zsb}.singleCompC 'docker ps --format \"{{.Names}}\"'" dkssh

