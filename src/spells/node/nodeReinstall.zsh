${zsb}.nodeReinstall() {
  local packageManager=$1
  local packageManagerLock=$2

  zparseopts -D -E -F -- l=keepLockFile n=keepNodeModules i=skipInstall c=cacheClean || return 1

  if [[ -z "$keepLockFile" && -f ${packageManagerLock} ]]
    then messageStatus rm ${packageManagerLock}
  fi

  if [[ -z "$keepNodeModules" && -d node_modules ]]
    then spinner rm -rf node_modules
  fi

  if [[ -n "$cacheClean" ]]
    then npm cache clean --force
  fi

  if [[ -n "$skipInstall" ]]
    then ${zsb}.cancel "Installation cancelled"
  fi

  ${packageManager} install && ${zsb}.isGitRepo && ${zsb}.gitStatus

  alert "🏁 ${packageManager} install 🏁"
}

compdef "_${zsb}.nonRepeatedListD '-l:Keep lock files' '-n:Keep node modules' '-i:Skip installation' '-c:Cache clean'" ${zsb}.nodeReinstall

alias npmreinstall="${zsb}.nodeReinstall npm package-lock.json"

alias yarnreinstall="${zsb}.nodeReinstall yarn yarn.lock"
