_tm.lastActive() tmux display-message -p '#S'

tm() {
  zparseopts -D -E -F -- c:=changeDir || return 1

  local session=$1
  if [[ -z "$session" ]]; then
    session=`_tm.lastActive`
    : ${session:='main'}
  fi

  # if outside tmux
  if [[ -z "$TMUX" ]]; then
    if tmux has-session -t $session 2>/dev/null
      then tmux attach-session -t $session
      else tmux new -s $session $changeDir
    fi

    return $?
  fi

  [[ "`_tm.lastActive`" = "$session" ]] && \
    ${zsb}.throw 'You did not move.'

  # Create session if it doesn't exists
  tmux new -s $session $changeDir -d 2>/dev/null

  tmux switch-client -t $session
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

