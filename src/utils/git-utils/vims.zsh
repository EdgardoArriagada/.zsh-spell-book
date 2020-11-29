vims() {
  vim "$@" && ${zsb}.isGitRepo && ${zsb}.gitStatus
}
