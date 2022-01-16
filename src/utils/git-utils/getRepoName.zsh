getRepoName() {
  local pathToDotGit=`git rev-parse --git-dir 2>/dev/null`
  [[ -z "$pathToDotGit" ]] && return 0

  if [[ "$pathToDotGit" = ".git" ]]; then
    printf ${PWD##*/}
  else
    local parentRootDir=`dirname ${pathToDotGit}`
    printf ${parentRootDir##*/}
  fi
}
