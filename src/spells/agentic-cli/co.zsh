co() {
  case "$#:${1-}" in
    0:)
      codex --dangerously-bypass-approvals-and-sandbox
      ;;
    1:--update-pr-title-and-description)
      codex exec 'Update this pr title and description'
      ;;
    1:--commit)
      codex exec 'Create a Git commit from the currently staged files'
      ;;
    *:--code-review)
      codex "\$thermo-nuclear-code-quality-review ${@:2}"
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
