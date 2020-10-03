tm() {
  tmux attach-session -t $1>/dev/null 2>&1 || tmux new -s $1>/dev/null 2>&1
}

# If any, complete with a list of current tmux sessions
complete -C "tmux ls 2>&1 | cut -d':' -s -f1" tm
