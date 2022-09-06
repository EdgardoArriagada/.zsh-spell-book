tmheaderbar() {
  tmux split-window -vb -p 6 -t 0 \; \
    send-keys -t 0 " printCentered \"${*}\" | less; exit " Enter \; \
    select-pane -t 1
}
