tmain() {
  local sessionName="main"

  if tmux has-session -t "$sessionName" 2>/dev/null; then
      tm "$sessionName"
      return 1
  fi

  tmux new-session -s ${sessionName} \; \
    new-window \; \
    split-window -hd -p 25 -t ${sessionName}:2.0 \; \
    split-window -vd -t ${sessionName}:2.1 \; \
    clock -t ${sessionName}:2.2 \; \
    send-keys -t ${sessionName}:2.1 ' watch sensors' Enter \; \
    select-pane -t ${sessionName}:2.0 \; \
    attach
}

