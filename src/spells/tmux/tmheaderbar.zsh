tmheaderbar() {
  local winId=$1
  shift 1

  tmux split-window -vb -p 6 -t ${winId}.0 -P -F "#{pane_id}" \; \
    send-keys -t ${winId}.0 " printCentered \"${*}\" | less; exit " Enter \; \
    select-pane -t ${winId}.1
}
