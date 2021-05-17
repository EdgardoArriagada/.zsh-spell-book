tml() {
  local latestSession=$(tmux display-message -p '#S')
  local currentList=(`eval "tmls"`)
  echo "\n"
  for session in "${currentList[@]}"; do
    if [[ "$session" = "$latestSession" ]]; then
      [[ -z "$TMUX" ]] && echo "~ ${session}" || echo "* $(hl ${session})"
    else
      echo "  $session"
    fi
  done
  echo "\n"
}

