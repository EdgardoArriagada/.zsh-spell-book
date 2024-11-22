tm() {
  zparseopts -D -E -F -- c:=changeDir || return 1

  # if outside tmux
  if [[ -z "$TMUX" ]]; then
    local session=${1:=`getLastActiveSessionName`}
    : ${session:=main}

    tmux new -A -s $session $changeDir
    return $?
  fi

  local session=${1:?Error: session name is required.}
  if [[ "$session" = "`getLastActiveSessionName`" ]]
    then ${zsb}.throw 'You did not move.'
  fi

  # Create session if it doesn't exists
  tmux new -s $session $changeDir -d 2>/dev/null

  tmux switch-client -t $session

  local winCount=`tmux display-message -p '#{session_windows}'`
  (( $winCount > 1 )) && tmux kill-window
}

_${zsb}.tm() {
  (( $CURRENT > 2 )) && return 0

  local tmuxList=(`tmls`)

  # if outside tmux
  [[ -z "$TMUX" ]] &&  _describe 'command' tmuxList && return 0

  local currentSession=(`getLastActiveSessionName`)

  tmuxList=(${tmuxList:|currentSession})

  _describe 'command' tmuxList
}

hisIgnore tm

compdef _${zsb}.tm tm

alias TM="toggleCapsLock && tm"

