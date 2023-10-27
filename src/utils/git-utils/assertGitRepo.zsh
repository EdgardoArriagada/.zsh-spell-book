${zsb}.assertGitRepo() {
  ${zsb}.isGitRepo || ${zsb}.throw "Not in a git repo."
}
