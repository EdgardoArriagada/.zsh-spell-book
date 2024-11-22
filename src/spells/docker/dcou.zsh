dcou() {
  if (( $# > 0 )); then
    printAndRun "docker-compose up $@"
    return $?
  fi

  # export it so it can be used in shutDown
  export winId=`${zsb}.tmux.getWinId`

  function shutDown() {
    hl_tmux_tab $winId --clear
  }

  TRAPINT() { shutDown }
  TRAPEXIT() { shutDown }

  hl_tmux_tab $winId

  printAndRun "docker-compose up"
}

alias dcoud="printAndRun 'docker-compose up -d'"
