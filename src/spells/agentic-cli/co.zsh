codex() {
  if [[ ${1-} == resume && ${2-} =~ '^[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$' && -n ${TMUX_PANE-} ]]; then
    zsb_tmux_agent_notification --bind-codex "$2" "$TMUX_PANE" || print -u2 'codex: tmux session binding failed'
  fi
  command codex "$@"
}

co() {
  case "$#:${1-}" in
    0:)
      codex --model gpt-6-sol --dangerously-bypass-approvals-and-sandbox
      ;;
    1:--update-pr-title-and-description)
      codex exec --dangerously-bypass-approvals-and-sandbox "Update this pr title and description.
- Use gh cli
- Don't wait for confirmation, just do it"
      ;;
    1:--commit)
      codex exec --dangerously-bypass-approvals-and-sandbox 'Create a Git commit from the currently staged files'
      ;;
    *:--code-review)
      codex --dangerously-bypass-approvals-and-sandbox "\$thermo-nuclear-code-quality-review ${@:2}"
      ;;
    *)
      codex "$@"
      ;;
  esac
}

_${zsb}.co() {
  (( CURRENT > 2 )) && return 0

  local -a options=(
    '--update-pr-title-and-description:update PR title and description'
    '--code-review:run thermo-nuclear code quality review'
    '--commit:create a git commit'
    'resume:resume a session'
    'fork:fork a session'
  )
  _describe 'option' options
}

compdef _${zsb}.co co
