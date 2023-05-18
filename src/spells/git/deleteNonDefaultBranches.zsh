deleteNonDefaultBranches() (
  local this=$0
  local nonDefaultBranches

  ${this}.setNonDefaultBranches() {
    nonDefaultBranches=`${zsb}.gitBranches | grep -vE "^(${ZSB_GIT_DEFAULT_BRANCHES})$"`
  }

  ${this}.printPrompt() {
    ${zsb}.confirmMenu.withItems \
      "The following branches will be deleted:" \
      "$nonDefaultBranches"
  }

  ${this}.performDeletion() {
    git branch -D <<< ${nonDefaultBranches}
  }

  { # main
    ${zsb}.validateGitRepo

    ${zsb}.userWorkingOnDefaultBranch || ${zsb}.throw "You must run this command from a default branch"

    ${this}.setNonDefaultBranches

    [[ -z "$nonDefaultBranches" ]] &&
      ${zsb}.cancel "There are no non default branches to delete."

    ${this}.printPrompt
    ${this}.performDeletion
  }
)

hisIgnore deleteNonDefaultBranches

_${zsb}.nocompletion deleteNonDefaultBranches

