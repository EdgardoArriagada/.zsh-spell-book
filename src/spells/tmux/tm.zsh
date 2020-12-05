tm() {
  local targetSession="$1"
  local DEFAULT_SESSION="main"

  [ -z "$targetSession" ] && targetSession="$DEFAULT_SESSION"

  if [ ! -z "$TMUX" ]; then # if inside tmux session
    if tmux has-session -t "$targetSession" 2>/dev/null; then
      echo "${ZSB_ERROR} Can't connect to $(hl ${targetSession}) whitin a tmux session."
      return 1
    fi
    tmux new-session -d -s "$targetSession" &&
      echo "${ZSB_SUCCESS} $(hl ${targetSession}) created."
    return $?
  fi

  tmux attach-session -t "$targetSession">/dev/null 2>&1 || \
    tmux new -s "$targetSession">/dev/null 2>&1
}

# If any, complete with a list of current tmux sessions
complete -C "tmux ls 2>&1 | cut -d':' -s -f1" tm
