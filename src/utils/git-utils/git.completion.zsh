# All completions must work from current working dir (not only from project root)
# Example usage: complete -C $ZSB_GIT_STAGED_GILES mycmd

${zsb}_getGitStagedFiles(){
  ${zsb}_isGitRepo && git status --short | grep '^[MARCD]' | sed s/^...//
}

${zsb}_getGitUnstagedFiles() {
  ${zsb}_isGitRepo && git status --short | grep -E '(^ [MARCD])|(^[MARCD]{2})' | sed s/^...//
}

${zsb}_getGitUntrackedFiles() {
  ${zsb}_isGitRepo && git status --short | grep -E '(^ \?)|(^\?{2})' | sed s/^...//
}

${zsb}_getGitUnstagedAndUntrackedFiles() {
  ${zsb}_isGitRepo && git status --short | grep -E '(^ [MARCD\?])|(^[MARCD\?]{2})' | sed s/^...//
}

${zsb}_getGitModifiedFiles() {
  ${zsb}_isGitRepo && git status --short | sed s/^...// || ls
}

${zsb}_gitBranches() {
  ${zsb}_isGitRepo && git branch | sed 's/^\*//'
}
