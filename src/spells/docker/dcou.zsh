dcou() {
  if (( $# > 0 )); then
    printAndRun "docker-compose up $@"
    return $?
  fi

  # export it so it can be used in shutDown
  export target=`${zsb}.tmux.getTarget`

  function shutDown() {
    hl_tmux_tab $target --clear
  }

  TRAPINT() { shutDown }
  TRAPEXIT() { shutDown }

  hl_tmux_tab $target

  printAndRun "docker-compose up"
}

alias dcoud="printAndRun 'docker-compose up -d'"
