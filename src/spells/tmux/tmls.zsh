# zsh completions does not read aliases, so tmls is a function
tmls() { tmux ls 2>&1 | cut -d':' -s -f1 }

tml() {
  local latestSession=`tmux display-message -p '#S'`
  local currentList=(`tmls`)
  printf "\n"
  for session in "${currentList[@]}"; do
    if [[ "$session" = "$latestSession" ]]; then
      [[ -z "$TMUX" ]] && printf "~ ${session}\n" || printf "* $(hl ${session})\n"
    else
      printf "  ${session}\n"
    fi
  done
  printf "\n"
}
