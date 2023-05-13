# zsh completions does not read aliases, so tmls is a function
tmls() tmux ls -F "#{session_name}"

tml() {
  local latestSession=`tmux display-message -p '#S'`
  local currentList=(`tmls`)

  <<< ''

  for session in "${currentList[@]}"; do
    if [[ "$session" = "$latestSession" ]]
      then [[ -z "$TMUX" ]] && <<< "~ ${session}" || <<< "* `hl ${session}`"
      else <<< "  ${session}"
    fi
  done

  <<< ''
}
