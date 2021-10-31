galock() (
  local lockFiles=(
    package-lock.json
    Gemfile.lock
    yarn.lock
    Cargo.lock
  )

  local redGitFiles=( $(${zsb}.getGitFiles 'red-safe') )
  local lockFilesToAdd=( ${lockFiles:*redGitFiles} )

  [[ -z "$lockFilesToAdd" ]] && return 0

  git add "${(z)lockFilesToAdd}" &&
    ${zsb}.info "$(hl "${lockFilesToAdd[@]}") added."
)

_${zsb}.nocompletion galock

