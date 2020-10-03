isGitRepo() {
  local pathToDotGit=$(git rev-parse --git-dir 2>/dev/null)
  return $([ ! -z "$pathToDotGit" ])
}
