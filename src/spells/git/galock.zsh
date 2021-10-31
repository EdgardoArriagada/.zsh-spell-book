galock() (
  local redGitFiles=( $(${zsb}.getGitFiles 'red-safe') )
  local lockFilesToAdd=( ${ZSB_GIT_LOCK_FILES:*redGitFiles} )

  [[ -z "$lockFilesToAdd" ]] && return 0

  git add "${(z)lockFilesToAdd}" &&
    ${zsb}.info "$(hl "${lockFilesToAdd[@]}") added."
)

_${zsb}.nocompletion galock

