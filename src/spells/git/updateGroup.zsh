${zsb}.updateGroup() (
  local this="$0"
  local repoFile="$1"
  local pass
  local repoList=($(< ${repoFile}))
  local currentRepo
  local currentBranch

  ${this}.main() {
    ${this}.validateRepoList

    ${this}.inputCredentials
    ${this}.manageEachRepo
    ${zsb}.success "Finished"
  }

  ${this}.validateRepoList() {
    [[ -z "$repoList" ]] && ${zsb}.throw "Invalid repolist."
  }

  ${this}.inputCredentials() {
    ${zsb}.prompt "Enter Credentials:"
    read -s pass
  }

  ${this}.manageEachRepo() {
    for repo in "${repoList[@]}"; do
      ${zsb}.fillWithToken '_'

      currentRepo="$repo"
      ${this}.validateRepo

      ${this}.setShellState

      if ${this}.isRepoInCleanState; then
        ${this}.printCleanHeader
        ${this}.updateRepo
      else
        ${this}.printDirtyHeader
        ${this}.manageManually
      fi
    done
  }

  ${this}.validateRepo() {
    [[ -d "$HOME/$currentRepo" ]] && return 0

    ${zsb}.throw "$(hl $currentRepo) does not exists, aborting.\nPlease check $(hl $repoFile)"
  }

  ${this}.setShellState() {
    builtin cd $HOME/$currentRepo
    currentBranch="$(git branch --show-current)"
  }

  ${this}.isRepoInCleanState() {
    local gitStatusOutput=$(script -qc "git status --short" /dev/null < /dev/null)
    ${zsb}.userWorkingOnDefaultBranch && [[ -z "$gitStatusOutput" ]]
  }

  ${this}.printCleanHeader() {
    printf "[ $(hl $currentRepo) ($currentBranch) ]\n"
  }

  ${this}.printDirtyHeader() {
    printf "[ $(hl $currentRepo) ($currentBranch) ± ]\n"
  }

  ${this}.updateRepo() {
    local withCredUrl="$(${this}.getWithCredUrl)"

    ${zsb}.info "Pulling from ${currentBranch}"
    git pull "$withCredUrl" "$currentBranch"
  }

  ${this}.getWithCredUrl() {
    local originUrl="$(git config --get remote.origin.url)"
    local regexReplacer="s/@/:${pass}@/"
    printf "$originUrl" | sed -E "$regexReplacer"
  }

  ${this}.manageManually() {
    ${zsb}.prompt "Enter $(hl n) to continue:"

    while true; do
      read key
      case $key in
      [Nn]*)
        # continue
        return 0 ;;
      esac
    done
  }

  ${this}.main "$@"
)
