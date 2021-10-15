${zsb}.updateGroup() (
  local this="$0"
  local repoFile="$1"
  local inputPass
  local repoList=($(< ${repoFile}))
  local currentRepo
  local currentBranch

  ${this}.main() {
    ${this}.validateRepoList

    ${this}.setCredentialsWithPrompt

    ${this}.manageEachRepo

    ${zsb}.success "Finished"
  }

  ${this}.validateRepoList() {
    [[ -z "$repoList" ]] && ${zsb}.throw "Invalid repolist."
  }

  ${this}.setCredentialsWithPrompt() {
    ${zsb}.prompt "Enter Credentials:"
    read -s inputPass
  }

  ${this}.manageEachRepo() {
    for repo in "${repoList[@]}"; do
      ${zsb}.fillWithToken '_'

      # Repo has to start with "~/"
      currentRepo="${HOME}/${repo#*/}"
      ${this}.validateRepo

      builtin cd "$currentRepo"
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

      if ${this}.isRepoInCleanState; then
        ${this}.printCleanHeader
        ${this}.updateRepo
        return 0
      fi
      ${this}.printDirtyHeader
      ${this}.playManageManuallyPrompt
      continueUpdating=$?
    done
  }

  ${this}.isRepoInCleanState() {
    local -r gitStatusOutput=$(script -qc "git status --short" /dev/null < /dev/null)
    ${zsb}.userWorkingOnDefaultBranch && [[ -z "$gitStatusOutput" ]]
  }

  ${this}.printHeader() {
    printf "[ $(hl ${currentRepo##*/}) ($currentBranch) ${1}]\n"
  }

  ${this}.printCleanHeader() ${this}.printHeader

  ${this}.printDirtyHeader() ${this}.printHeader "± "

  ${this}.updateRepo() {
    local -r withCredUrl="$(${this}.getWithCredUrl)"

    ${zsb}.info "Pulling from ${currentBranch}"
    git pull "$withCredUrl" "$currentBranch"
  }

  ${this}.getWithCredUrl() {
    local -r originUrl="$(git config --get remote.origin.url)"
    local -r regexReplacer="s/@/:${inputPass}@/"
    printf "$originUrl" | sed -E "$regexReplacer"
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

  ${this}.main "$@"
)
