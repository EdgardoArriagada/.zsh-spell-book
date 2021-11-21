tmonly() {
  local -r tmuxList=( $(tmls) )
  local -r currentSession=( "$(tmux display-message -p '#S')" )
  local -r otherSessions=${tmuxList:|currentSession}

  [[ -z "$otherSessions" ]] && ${zsb}.cancel "There are no other sessions."

  ${zsb}.confirmMenu.withItems \
    "You are about to kill the following sessions" \
    "$otherSessions"

  tmux kill-session -a
}

_${zsb}.nocompletion tmonly
