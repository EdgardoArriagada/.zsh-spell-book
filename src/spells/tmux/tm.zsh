.tm.getCurrenSession() tmux display-message -p '#S'

tm() {
  local -r targetSession=${1:='main'}

  # if outside tmux session
  [[ -z "$TMUX" ]] && tmux new -A -s "$targetSession" && return $?

  [[ "`.tm.getCurrenSession`" = "$targetSession" ]] && \
    ${zsb}.throw "You did not move."

  # Create session if it doesn't exists
  tmux new -d -s "$targetSession" 2>/dev/null

  tmux switch-client -t "$targetSession"
}

_${zsb}.tm() {
  [[ "$CURRENT" -gt "2" ]] && return 0

  local -r tmuxList=( `tmls` )
  if [[ -n "$TMUX" ]]; then # if inside tmux session
    local -r currentSession=( `.tm.getCurrenSession` )
    local -r actualList=(${tmuxList:|currentSession})
    _describe 'command' actualList
    return
  fi

  _describe 'command' tmuxList
}

hisIgnore tm

compdef _${zsb}.tm tm

alias TM="toggleCapsLock && tm"

