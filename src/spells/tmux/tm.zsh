tm() {
  local -r targetSession=${1:='main'}

  # if outside tmux session
  [[ -z "$TMUX" ]] && tmux new -A -s "$targetSession" && return $?

  local -r currentSession="$(tmux display-message -p '#S')"
  [[ "$currentSession" = "$targetSession" ]] && ${zsb}.throw "You did not move."

  if ! tmux has-session -t "$targetSession" 2>/dev/null; then
    tmux new -d -s "$targetSession"
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

alias TM="toggleCapsLock && tm"

