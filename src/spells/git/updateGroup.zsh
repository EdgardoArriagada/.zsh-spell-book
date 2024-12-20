${zsb}.updateGroup() (
  local this=$0
  local repoFile=$1
  local repoList=($(< $repoFile))
  local currentRepo
  local currentBranch

  ${this}.validateRepoList() {
    [[ -z "$repoList" ]] && ${zsb}.throw "Invalid repolist."
  }

  ${this}.manageEachRepo() {
    for repo in "${repoList[@]}"; do
      hr

      # Repo has to start with "~/"
      currentRepo=$HOME/${repo#*/}
      ${this}.validateRepo

      builtin cd $currentRepo
      ${this}.manageRepo
    done
  }

  ${this}.validateRepo() {
    [[ -d "$currentRepo" ]] && return 0

    ${zsb}.throw "$(hl $currentRepo) does not exists, aborting.\nPlease check $(hl $repoFile)"
  }

  ${this}.setCurrentBranch() currentBranch="$(git branch --show-current)"

  ${this}.manageRepo() {

    local continueUpdating=0
    while [[ "$continueUpdating" == "0" ]]; do
      ${this}.setCurrentBranch

      if ${this}.isEligibleForUpdate; then
        ${this}.printCleanHeader
        ${this}.updateRepo
        return 0
      fi
      ${this}.printDirtyHeader
      ${this}.playManageManuallyPrompt
      continueUpdating=$?
    done
  }

  ${this}.isEligibleForUpdate() {
    ${zsb}.isUserOnDefaultBranch && ${zsb}.isRepoClean
  }

  ${this}.printHeader() <<< "[ $(hl `get_repo_name`) ($currentBranch) ${1}]"

  ${this}.printCleanHeader() ${this}.printHeader

  ${this}.printDirtyHeader() ${this}.printHeader "± "

  ${this}.updateRepo() {
    ${zsb}.info "Pulling from $currentBranch"
    git pull origin "$currentBranch"
  }

  ${this}.playManageManuallyPrompt() {
    ${zsb}.prompt "$(hl '[r,n]')"

    while true; do
      read key
      case $key in
      [Nn]*)
        # continue updating
        return 1 ;;
      [Rr]*)
        # retry updating
        return 0 ;;
      esac
    done
  }

  { # main
    ${this}.validateRepoList

    ${this}.manageEachRepo

    ${zsb}.success "Finished"
  }
)
