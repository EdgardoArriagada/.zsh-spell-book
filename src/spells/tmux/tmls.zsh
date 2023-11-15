# zsh completions does not read aliases, so tmls is a function
tmls() tmux ls -F "#{session_name}" 2>/dev/null

tml() {
  local latestSession=`tmux display-message -p '#S' 2>/dev/null`
  <<< ''

  for session in `tmls`; do
    if [[ "$session" = "$latestSession" ]]
      then [[ -z "$TMUX" ]] && <<< "~ $session" || <<< "* `hl $session`"
      else <<< "  $session"
    fi
  done

  <<< ''
}
