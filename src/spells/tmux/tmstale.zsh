tmstale() {
  local session=$(tmdiff | fzf --exit-0 --select-1)
  [[ -z "$session" ]] && ${zsb}.cancel 'There are no stale sessions.'

  if [[ -n "$TMUX" ]]
    then tmux switch-client -t "=$session"
    else tmux attach-session -t "=$session"
  fi
}

hisIgnore tmstale
