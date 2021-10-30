deleteNonDefaultBranches() (
  local this="$0"
  local nonDefaultBranches

  ${this}.main() {
    ${zsb}.isGitRepo || ${zsb}.throw "You must run this command inside a git project."

    ${zsb}.userWorkingOnDefaultBranch || ${zsb}.throw "You must run this command from a default branch"

    ${this}.setNonDefaultBranches

    [[ -z "$nonDefaultBranches" ]]
      ${zsb}.cancel "There are no non default branches to delete."

    ${this}.printPrompt
    ${this}.playOptionsMenu
  }

  ${this}.setNonDefaultBranches() {
    nonDefaultBranches="$(git branch --no-color | command grep -vE "^(\+|\*|\s*${ZSB_GIT_DEFAULT_BRANCHES}\s*$)")"
  }

  ${this}.printPrompt() {
    ${zsb}.fullPrompt \
      "The following branches will be deleted:" \
      "$nonDefaultBranches"
  }

  ${this}.playOptionsMenu() {
    ${zsb}.confirmMenu && ${this}.performDeletion
  }

  ${this}.performDeletion() {
    echo "$nonDefaultBranches" | xargs git branch -D
  }

  ${this}.main "$@"
)

_${zsb}.nocompletion deleteNonDefaultBranches

