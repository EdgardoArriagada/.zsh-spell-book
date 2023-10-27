dcossh() {
  local inputService=${1?Error: Please provide a docker services.}
  printAndRun "docker-compose exec ${inputService} bash"
}

compdef "_${zsb}.cachedSingleCompWD 'docker-compose ps --services' dcossh 30" dcossh

