tmonly() {
  local tmuxList=( $(tmls) )
  local currentSession=( "$(tmux display-message -p '#S')" )
  local otherSessions=${tmuxList:|currentSession}

  [[ -z "$otherSessions" ]] && ${zsb}.cancel "There are no other sessions."

  ${zsb}.confirmMenu.withItems \
    "You are about to kill the following sessions" \
    "$otherSessions"

  tmux kill-session -a
}

_${zsb}.nocompletion tmonly
