# Automatically add lock files
galock() (
  declare -A lockFiles
  lockFiles[package-lock.json]=true
  lockFiles[Gemfile.lock]=true
  lockFiles[yarn.lock]=true

  local redGitFiles=( $(${zsb}.getGitFiles 'red') )
  local hasAFileBeenAdded=false

  for file in "${redGitFiles[@]}"; do
    if [ "${lockFiles[$file]}" = "true" ]; then
      git add "$file"
      hasAFileBeenAdded=true
    fi
  done

  if [ "$hasAFileBeenAdded" = "false" ]; then
    echo "${ZSB_WARNING} No $(hl "*.lock") files have been added."
  fi

  ${zsb}.gitStatus
)

complete galock

