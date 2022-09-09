tmheaderbar() {
  tmux split-window -vb -p 6 -t 0 \; \
    send-keys -t 0 " printCentered \"${*}\" | less; exit " Enter \; \
    select-pane -t 1
}

# param $1: sessionName=`tmux display-message -p '#S'`
# param $2: windowNumber=`tmux display-message -p '#I'`
tmheaderbarFixed() {
  local givenWindow="${1}:${2}"
  shift 2

  tmux split-window -vb -p 6 -t ${givenWindow}.0 \; \
    send-keys -t ${givenWindow}.0 " printCentered \"${*}\" | less; exit " Enter \; \
    select-pane -t ${givenWindow}.1

}
