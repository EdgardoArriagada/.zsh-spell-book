tm() {
  local targetSession="$1"
  local DEFAULT_SESSION="main"

  [ -z "$targetSession" ] && targetSession="$DEFAULT_SESSION"

  tmux attach-session -t "$targetSession">/dev/null 2>&1 || \
    tmux new -s "$targetSession">/dev/null 2>&1
}

# If any, complete with a list of current tmux sessions
complete -C "tmux ls 2>&1 | cut -d':' -s -f1" tm
