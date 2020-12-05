tmprod() {
  local sessionName="PROD"

  if tmux has-session -t "$sessionName" 2>/dev/null; then
      tm "$sessionName"
      return 1
  fi

  tmux new-session -s ${sessionName} \; \
    split-window -hd -p 80 \; \
    clock -t ${sessionName}:0.0 \; \
    select-pane -t ${sessionName}:0.1 \; \
    attach
}

