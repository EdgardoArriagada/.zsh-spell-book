galock() (
  local lockFiles=(
    package-lock.json
    Gemfile.lock
    yarn.lock
    Cargo.lock
  )

  local redGitFiles=( $(${zsb}.getGitFiles 'red-safe') )
  local lockFilesToAdd=( ${lockFiles:*redGitFiles} )

  if [[ -z "$lockFilesToAdd" ]]; then
    ${zsb}.warning "No $(hl "*.lock") files have been added."
  else
    git add "${(z)lockFilesToAdd}"
  fi

  ${zsb}.gitStatus
)

_${zsb}.nocompletion galock

