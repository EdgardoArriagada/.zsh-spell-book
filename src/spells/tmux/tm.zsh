tm() {
  local targetSession="$1"
  local DEFAULT_SESSION="main"

  [ -z "$targetSession" ] && targetSession="$DEFAULT_SESSION"

  if [ ! -z "$TMUX" ]; then # if inside tmux session
    if ! tmux has-session -t "$targetSession" 2>/dev/null; then
      tmux new-session -d -s "$targetSession"
    fi
    tmux switch-client -t "$targetSession"
  fi

  tmux attach-session -t "$targetSession">/dev/null 2>&1 || \
    tmux new -s "$targetSession">/dev/null 2>&1
}

_${zsb}.tm() {
  [ "$COMP_CWORD" -gt "1" ] && return 0

  COMPREPLY=( $(tmls) )
}

complete -F _${zsb}.tm tm
