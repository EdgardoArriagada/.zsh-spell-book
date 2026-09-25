#!/usr/bin/env zsh
set -euo pipefail
typeset -r charms_dir=${0:A:h:h:h}/charms

main() {
  local id=123e4567-e89b-12d3-a456-426614174000
  local capture_line="header
│  Session:                $id                       │
footer
footer"

  tmux() {
    case $1 in
      display-message) print -r -- 100 ;;
      show-window-options) print -r -- manual ;;
      capture-pane) print -r -- "$capture_line" ;;
      new-window) print -r -- "fork:$5:$6" ;;
    esac
  }
  zsb_tmux_agent_notification() { print -r -- "bind:$2:$3" }

  [[ $(source "$charms_dir/zsb_charm_codex_session_id" %9) == "$id" ]]
  capture_line=$(printf ' · %s · \nfooter' "$id")
  [[ $(source "$charms_dir/zsb_charm_codex_session_id" %9) == "$id" ]]
  capture_line='no session'
  if ( source "$charms_dir/zsb_charm_codex_session_id" %9 ) >/dev/null; then return 1; fi
  zsb_charm_codex_session_id() { print -r -- "$id" }

  [[ $(source "$charms_dir/zsb_charm_codex_bind" %9) == "bind:$id:%9" ]]

  [[ $(source "$charms_dir/zsb_charm_codex_fork" %9 /tmp/project) == "fork:/tmp/project:codex fork $id" ]]
}

main
