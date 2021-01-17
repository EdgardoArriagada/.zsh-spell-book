tmsidebar() {
  tmux split-window -hd -p 25 -t 0 \; \
    split-window -vd -t 1 \; \
    clock -t 2 \; \
    send-keys -t 1 ' watch sensors' Enter \; \
    select-pane -t 0
}
