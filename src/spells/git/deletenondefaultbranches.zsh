deletenondefaultbranches() (
  local nonDefaultBranches

  main() {
    if ! ${zsb}_isGitRepo; then
      echo "${ZSB_ERROR} You must run this command inside a git project"
      return 1
    fi

    if ! userWorkingOnDefaultBranch; then
      echo "${ZSB_ERROR} You must run this command from a default branch"
      return 1
    fi

    setNonDefaultBranches

    if [ -z "$nonDefaultBranches" ]; then
      echo "${ZSB_INFO} There are no non default branches to delete."
      return 0
    fi

    printPrompt
    playOptionsMenu
  }

  setNonDefaultBranches() {
    nonDefaultBranches="$(git branch --no-color | command grep -vE "^(\+|\*|\s*${ZSB_GIT_DEFAULT_BRANCHES}\s*$)")"
  }

  printPrompt() {
    echo "${ZSB_WARNING} The following branches will be deleted:"
    echo " "
    echo "$(hl "$nonDefaultBranches")"
    echo " "
    echo "${ZSB_PROMPT} Do you really want to delete these branches? [Y/n]"
  }

  playOptionsMenu() {
    ${zsb}_yesNoMenu performDeletion &&
      echo "$ZSB_SUCCESS: non default branches deleted"
  }

  performDeletion() {
    echo "$nonDefaultBranches" | xargs git branch -D
  }

  main "$@"
)
