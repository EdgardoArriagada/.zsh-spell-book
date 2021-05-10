vims() {
 eval "nnvim $@" && ${zsb}.isGitRepo && ${zsb}.gitStatus
}
