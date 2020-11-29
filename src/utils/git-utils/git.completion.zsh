# All completions must work from current working dir (not only from project root)
# Example usage: complete -C $ZSB_GIT_STAGED_GILES mycmd

${zsb}.getGitStagedFiles(){
  ${zsb}.isGitRepo && git status --short | grep '^[MARCD]' | sed s/^...//
}

${zsb}.getGitUnstagedFiles() {
  ${zsb}.isGitRepo && git status --short | grep -E '(^ [MARCD])|(^[MARCD]{2})' | sed s/^...//
}

${zsb}.getGitUntrackedFiles() {
  ${zsb}.isGitRepo && git status --short | grep -E '(^ \?)|(^\?{2})' | sed s/^...//
}

${zsb}.getGitUnstagedAndUntrackedFiles() {
  ${zsb}.isGitRepo && git status --short | grep -E '(^ [MARCD\?])|(^[MARCD\?]{2})' | sed s/^...//
}

${zsb}.getGitModifiedFiles() {
  ${zsb}.isGitRepo && git status --short | sed s/^...// || ls
}

${zsb}.gitBranches() {
  ${zsb}.isGitRepo && git branch | sed 's/^\*//'
}
