_tm.lastActive() tmux display-message -p '#S' 2>/dev/null

tm() {
  zparseopts -D -E -F -- c:=changeDir || return 1

  local lastActive=`_tm.lastActive`

  local session=${1:=$lastActive}
  : ${session:='main'}

  if [[ -n "$TMUX" ]] && [[ "$lastActive" = "$session" ]]
    then ${zsb}.throw 'You did not move.'
  fi

  # Create session if it doesn't exists
  tmux new -s $session $changeDir -d 2>/dev/null

  if [[ -n "$TMUX" ]]
    then tmux switch-client -t $session
    else tmux attach-session -t $session
  fi
}

_${zsb}.tm() {
  (( $CURRENT > 2 )) && return 0

  local tmuxList=( `tmls` )

  # if outside tmux
  [[ -z "$TMUX" ]] &&  _describe 'command' tmuxList && return 0

  local currentSession=( `_tm.lastActive` )
  local actualList=(${tmuxList:|currentSession})
  _describe 'command' actualList
}

hisIgnore tm

compdef _${zsb}.tm tm

alias TM="toggleCapsLock && tm"

