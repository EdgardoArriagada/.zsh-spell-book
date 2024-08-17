dcou() {
  if (( $# > 0 )); then
    printAndRun "docker-compose up $@"
    return $?
  fi

  function shutDown() {
    hl_tmux_tab --clear
  }

  TRAPINT() { shutDown }
  TRAPEXIT() { shutDown }

  hl_tmux_tab

  printAndRun "docker-compose up"
}

alias dcoud="printAndRun 'docker-compose up -d'"
