vims() {
  nvim "$@" && ${zsb}.isGitRepo && ${zsb}.gitStatus
}
