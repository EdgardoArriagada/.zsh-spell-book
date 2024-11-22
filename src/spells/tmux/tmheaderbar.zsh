tmheaderbar() {
  tmux split-window -vb -p 6 -t 0 -P -F "#{pane_id}" \; \
    send-keys -t 0 " printCentered \"${*}\" | less; exit " Enter \; \
    select-pane -t 1
}

# param $1: sessionName=`tmux display-message -p '#S'`
# param $2: windowNumber=`tmux display-message -p '#I'`
tmheaderbarFixed() {
  local winId=$1
  shift 1

  tmux split-window -vb -p 6 -t ${winId}.0 -P -F "#{pane_id}" \; \
    send-keys -t ${winId}.0 " printCentered \"${*}\" | less; exit " Enter \; \
    select-pane -t ${winId}.1
}
