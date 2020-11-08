vims() {
  vim "$@" && ${zsb}_isGitRepo && ${zsb}_gitStatus
}
