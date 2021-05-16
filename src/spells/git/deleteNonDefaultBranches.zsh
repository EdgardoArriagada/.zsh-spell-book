deleteNonDefaultBranches() (
  local this="$0"
  local nonDefaultBranches

  ${this}.main() {
    ${zsb}.isGitRepo || ${zsb}.throw "You must run this command inside a git project."

    ${zsb}.userWorkingOnDefaultBranch || ${zsb}.throw "You must run this command from a default branch"

    ${this}.setNonDefaultBranches

    if [[ -z "$nonDefaultBranches" ]]; then
      ${zsb}.info "There are no non default branches to delete."
      return 0
    fi

    ${this}.printPrompt
    ${this}.playOptionsMenu
  }

  ${this}.setNonDefaultBranches() {
    nonDefaultBranches="$(git branch --no-color | command grep -vE "^(\+|\*|\s*${ZSB_GIT_DEFAULT_BRANCHES}\s*$)")"
  }

  ${this}.printPrompt() {
    ${zsb}.warning "The following branches will be deleted:"
    echo " "
    echo "$(hl "$nonDefaultBranches")"
    echo " "
    ${zsb}.prompt "Do you really want to delete these branches? [Y/n]"
  }

  ${this}.playOptionsMenu() {
    ${zsb}.confirmMenu && ${this}.performDeletion &&
      ${zsb}.success "Non default branches deleted."
  }

  ${this}.performDeletion() {
    echo "$nonDefaultBranches" | xargs git branch -D
  }

  ${this}.main "$@"
)

_${zsb}.nocompletion deleteNonDefaultBranches

