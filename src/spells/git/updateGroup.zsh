${zsb}.updateGroup() (
  local this="$0"
  local repoList="$1"
  local pass
  local repoList=($(< ${repoList}))

  ${this}.main() {
    ${this}.validateRepoList

    ${this}.inputCredentials
    ${this}.manageEachRepo
    ${zsb}.success "Finished"
  }

  ${this}.validateRepoList() {
      [ -z "$repoList" ] && ${zsb}.throw "Invalid repolist."
  }

  ${this}.inputCredentials() {
    ${zsb}.prompt "Enter Credentials:"
    read -s pass
  }

  ${this}.manageEachRepo() {
    for repo in "${repoList[@]}"; do
      ${zsb}.fillWithToken '_'
      ${this}.manageRepo "$repo"
    done
  }

  ${this}.manageRepo() {
    local repo="$1"
    ${zsb}.info "Managing $(hl $repo)."

    if [ ! -d "$HOME/$repo" ]; then
      ${zsb}.throw "Repo does not exists, aborting."
    fi

    builtin cd $HOME/$repo

    if ${this}.isRepoInCleanState; then
      ${this}.updateRepo "$repo"
    else
      ${this}.manageManually
    fi
  }

  ${this}.isRepoInCleanState() {
    local gitStatusOutput=$(script -qc "git status --short" /dev/null < /dev/null)
    ${zsb}.userWorkingOnDefaultBranch && [ -z "$gitStatusOutput" ]
  }

  ${this}.updateRepo() {
    local originUrl="$(git config --get remote.origin.url)"
    local regexReplacer="s/@/:${pass}@/"
    local withCredUrl="$(printf "$originUrl" | sed -E "$regexReplacer")"
    git pull "$withCredUrl"
  }

  ${this}.manageManually() {
    ${zsb}.info "Stopped at dirty repo.\nOpen a new terminal to manage"
    ${zsb}.prompt "Enter $(hl c) or $(hl n) to continue:"

    while true; do
      read key
      case $key in
      [CcNn]*)
        ${zsb}.info "Continuing..."
        return 0
        ;;
      esac
    done
  }

  ${this}.main "$@"
)
