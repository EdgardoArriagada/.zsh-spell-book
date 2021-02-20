dcorr() {
  local session=$(tmux display-message -p '#S')
  local window=$(tmux display-message -p '#I')
  local pane=$(tmux display-message -p '#P')

  eval "printAndRun 'docker-compose down;\
    tmux kill-pane -a -t ${session}:${window}.${pane};\
    docker-compose up'\
  "
}
