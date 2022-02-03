tmheaderbar() {
  tmux split-window -vb -p 6 -t 0 \; \
    send-keys -t 0 " printCentered \"$*\" | less " Enter \; \
    select-pane -t 1
}
