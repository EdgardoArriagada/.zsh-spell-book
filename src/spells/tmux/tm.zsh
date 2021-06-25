tm() {
  local -r targetSession=${1:='main'}

  if [[ -z "$TMUX" ]]; then # if outside tmux session
    tmux attach-session -t "$targetSession">/dev/null 2>&1 || \
      tmux new -s "$targetSession">/dev/null 2>&1
    return $?
  fi

  local -r currentSession="$(tmux display-message -p '#S')"
  [[ "$currentSession" = "$targetSession" ]] && ${zsb}.throw "You did not move."

  if ! tmux has-session -t="$targetSession" 2>/dev/null; then
    tmux new-session -d -s "$targetSession"
  fi

  tmux switch-client -t "$targetSession"
}

_${zsb}.tm() {
  [[ "$CURRENT" -gt "2" ]] && return 0

  local -r tmuxList=( $(tmls) )
  if [[ ! -z "$TMUX" ]]; then # if inside tmux session
    local -r currentSession=( "$(tmux display-message -p '#S')" )
    local -r actualList=(${tmuxList:|currentSession})
    _describe 'command' actualList
    return
  fi

  _describe 'command' tmuxList
}

compdef _${zsb}.tm tm

