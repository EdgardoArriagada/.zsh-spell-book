#!/usr/bin/env zsh
set -euo pipefail
typeset -r charms_dir=${0:A:h:h:h}/charms

main() {
  local id=123e4567-e89b-12d3-a456-426614174000
  local tries=0
  local capture_line="header
$id · model · directory
footer"

  tmux() {
    case $1 in
      capture-pane) print -r -- "$capture_line" ;;
      new-window) print -r -- "fork:$5:$6" ;;
      display-message)
        if [[ ${2-} == -pt ]]; then
          (( tries )) && print codex || print zsh
        else
          print notice
        fi
        ;;
      resize-window) return 1 ;;
    esac
  }
  zsb_tmux_agent_notification() { print -r -- "bind:$2:$3" }

  [[ $(source "$charms_dir/zsb_charm_codex_session_id" %9) == "$id" ]]
  capture_line=$(printf '│ %s · model\nfooter' "$id")
  [[ $(source "$charms_dir/zsb_charm_codex_session_id" %9) == "$id" ]]
  capture_line='no session'
  if ( source "$charms_dir/zsb_charm_codex_session_id" %9 ) >/dev/null; then return 1; fi
  [[ -z $(source "$charms_dir/zsb_charm_codex_session_id" %9 --quiet) ]]
  zsb_charm_codex_session_id() { print -r -- "$id" }
  sleep() { (( ++tries )) }

  [[ $(source "$charms_dir/zsb_charm_codex_bind" %9) == "bind:$id:%9" ]]
  [[ $(source "$charms_dir/zsb_charm_codex_bind" %9 --wait) == "bind:$id:%9" ]]

  [[ $(source "$charms_dir/zsb_charm_codex_fork" %9 /tmp/project) == "fork:/tmp/project:codex fork $id" ]]
}

main
