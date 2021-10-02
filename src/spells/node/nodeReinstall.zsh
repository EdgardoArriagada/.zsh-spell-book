${zsb}.nodeReinstall() {
  local packageManager=$1
  local packageManagerLock=$2

  zparseopts -D -E -F -- l=keepLockFile n=keepNodeModules i=skipInstall || return 1

  if [[ -z "$keepLockFile" ]] && [[ -f ${packageManagerLock} ]]; then
    rm ${packageManagerLock} && ${zsb}.info "$(hl ${packageManagerLock}) deleted"
  fi

  if [[ -z "$keepNodeModules" ]] && [[ -d node_modules ]]; then
    rm -rf node_modules && ${zsb}.info "$(hl 'node_modules') deleted"
  fi

  if [[ -n "$skipInstall" ]]; then
    ${zsb}.info "Installation cancelled"; return 0
  fi

  # BUG: npm always tries to acess to zsb_lazy_load inside a bash function
  # eval is a workaround
  eval "${packageManager} install" && ${zsb}.isGitRepo && ${zsb}.gitStatus
}

compdef "_${zsb}.nonRepeatedListD '-l:Keep lock files' '-n:Keep node modules' '-i:Skip installation'" ${zsb}.nodeReinstall

alias npmreinstall="${zsb}.nodeReinstall npm package-lock.json"

alias yarnreinstall="${zsb}.nodeReinstall yarn yarn.lock"
