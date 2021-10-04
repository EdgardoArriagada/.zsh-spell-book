tmonly() {
  local -r tmuxList=( $(tmls) )
  local -r currentSession=( "$(tmux display-message -p '#S')" )
  local -r otherSessions=${tmuxList:|currentSession}

  [[ -z "$otherSessions" ]] && ${zsb}.info "There are no other sessions." && return 0

  ${zsb}.fullPrompt \
    "You are about to kill the following sessions" \
    "$otherSessions"

  ${zsb}.confirmMenu && tmux kill-session -a
}

_${zsb}.nocompletion tmonly
