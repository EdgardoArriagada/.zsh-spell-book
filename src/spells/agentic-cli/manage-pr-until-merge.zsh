manage-pr-until-merge() {
  emulate -L zsh
  local pr='' chat='' prUrl watcherPane dependency
  local -a selection watcherArgs prArgs chatArgs helpArgs
  local prPattern='^https://github\.com/[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+/pull/[1-9][0-9]*$'

  zparseopts -D -E -M -- p:=prArgs -pull-request:=p \
    c:=chatArgs -chat-integration:=c h=helpArgs -help=h || return 2
  watcherArgs=( "$@" )
  pr=${${prArgs[-1]-}#=}
  chat=${${chatArgs[-1]-}#=}
  if (( $#helpArgs )); then
    print -r -- 'Usage: manage-pr-until-merge [-p|--pull-request PR] [-c|--chat-integration HTTPS_URL] [watch-pr-events flags]
Run inside tmux: Codex in this pane (65%), watcher on the right (35%).
PR defaults to the current branch; accepts a PR number or GitHub PR URL.
Remaining flags go to watch-pr-events; PR and Codex thread are set automatically.'
    return 0
  fi

  if [[ -z ${TMUX-} || ${TMUX_PANE-} != %<-> ]]; then
    print -u2 -- 'manage-pr-until-merge: run inside a tmux pane'
    return 1
  fi
  if (( $#prArgs )); then
    if [[ ! $pr =~ '^[1-9][0-9]*$' && ! $pr =~ $prPattern ]]; then
      print -u2 -- 'manage-pr-until-merge: PR must be a number or GitHub PR URL'
      return 2
    fi
    selection=( "$pr" )
  fi
  if (( $#chatArgs )) && [[ ! $chat =~ '^https://[^[:space:]]+$' ]]; then
    print -u2 -- 'manage-pr-until-merge: chat integration must be an HTTPS URL'
    return 2
  fi
  for dependency in tmux gh codex watch-pr-events zsb_charm_find_codex_session_id; do
    if ! whence -p "$dependency" >/dev/null; then
      print -u2 -- "manage-pr-until-merge: $dependency not found"
      return 1
    fi
  done
  prUrl=$(gh pr view "${selection[@]}" --json url --jq .url) || return 1
  if [[ ! $prUrl =~ $prPattern ]]; then
    print -u2 -- 'manage-pr-until-merge: gh returned an invalid PR URL'
    return 1
  fi

  watcherPane=$(tmux split-window -h -d -l '35%' -t "$TMUX_PANE" -c "$PWD" -P -F '#{pane_id}' \
    zsh -f -c '
      pane=$1
      shift
      print -r -- "Waiting for Codex thread..."
      thread=$(zsb_charm_find_codex_session_id "$pane" 300) || {
        print -u2 -- "manage-pr-until-merge: timed out waiting for Codex thread"
        exit 1
      }
      exec watch-pr-events "$@" --codex-thread "$thread"
    ' manage-pr-until-merge "$TMUX_PANE" "${watcherArgs[@]}" --pull-request "$prUrl") || return 1

  {
    local prompt="\$manage-pr-until-merged $prUrl"
    [[ -z $chat ]] || prompt+=" $chat"
    codex --config 'tui.status_line=["session-id"]' "$prompt"
  } always {
    tmux kill-pane -t "$watcherPane" 2>/dev/null
  }
}
