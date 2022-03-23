${zsb}.nodeReinstall() {
  local packageManager=$1
  local packageManagerLock=$2

  zparseopts -D -E -F -- l=keepLockFile n=keepNodeModules i=skipInstall || return 1

  if [[ -z "$keepLockFile" ]] && [[ -f ${packageManagerLock} ]]; then
    rm ${packageManagerLock} && ${zsb}.info "$(hl ${packageManagerLock}) deleted"
  fi

  if [[ -z "$keepNodeModules" ]] && [[ -d node_modules ]]; then
    spinner rm -rf node_modules
  fi

  [[ -n "$skipInstall" ]] &&
    ${zsb}.cancel "Installation cancelled"

  # BUG: npm always tries to acess to ${zsb}.lazyLoad inside a bash function
  # eval is a workaround
  eval "${packageManager} install" && ${zsb}.isGitRepo && ${zsb}.gitStatus

  alert "DONE: '${packageManager} install'"
}

compdef "_${zsb}.nonRepeatedListD '-l:Keep lock files' '-n:Keep node modules' '-i:Skip installation'" ${zsb}.nodeReinstall

alias npmreinstall="${zsb}.nodeReinstall npm package-lock.json"

alias yarnreinstall="${zsb}.nodeReinstall yarn yarn.lock"
