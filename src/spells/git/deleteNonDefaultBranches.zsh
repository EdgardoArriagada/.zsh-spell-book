deleteNonDefaultBranches() (
  local this=$0
  local nonDefaultBranches

  ${this}.performDeletion() {
    print ${nonDefaultBranches} | xargs git branch -D
  }

  ${this}.printPrompt() {
    ${zsb}.confirmMenu.withItems \
      "The following branches will be deleted:" \
      "$nonDefaultBranches"
  }

  ${this}.assertNonDefaultBranches() {
    [[ -z "$nonDefaultBranches" ]] &&
      ${zsb}.cancel 'There are no non default branches to delete.'
  }

  ${this}.setNonDefaultBranches() {
    nonDefaultBranches=`${zsb}.gitBranches | grep -vE "^(${ZSB_GIT_DEFAULT_BRANCHES})$"`
    ${this}.assertNonDefaultBranches
  }

  ${this}.assertRepoConditions() {
    ${zsb}.assertGitRepo

    ${zsb}.userWorkingOnDefaultBranch ||
      ${zsb}.cancel 'You must run this command from a default branch'
  }

  { # main
    ${this}.assertRepoConditions

    ${this}.setNonDefaultBranches

    ${this}.printPrompt
    ${this}.performDeletion
  }
)

hisIgnore deleteNonDefaultBranches

_${zsb}.nocompletion deleteNonDefaultBranches

