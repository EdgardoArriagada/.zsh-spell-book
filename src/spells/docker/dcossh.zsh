dcossh() {
  local -r inputService=${1:?'Please provide a docker services'}

  printAndRun "docker-compose exec ${inputService} bash"
}

compdef "_${zsb}.cachedSingleCompC 'docker-compose ps --services' dcossh 30" dcossh

