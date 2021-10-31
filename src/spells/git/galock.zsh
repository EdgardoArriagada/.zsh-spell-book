galock() (
  local lockFiles=(
    package-lock.json
    Gemfile.lock
    yarn.lock
    Cargo.lock
  )

  local redGitFiles=( $(${zsb}.getGitFiles 'red-safe') )
  local lockFilesToAdd=( ${lockFiles:*redGitFiles} )

  [[ -z "$lockFilesToAdd" ]] &&
    ${zsb}.cancel "No $(hl "*.lock") files have been added."

  git add "${(z)lockFilesToAdd}"
  ${zsb}.gitStatus
)

_${zsb}.nocompletion galock

