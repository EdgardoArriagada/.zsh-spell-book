${zsb}.nodeReinstall() {
  local packageManager=$1
  local packageManagerLock=$2

  zparseopts -D -E -F -- l=keepLockFile n=keepNodeModules i=skipInstall f=cacheClean || return 1

  if [[ -n "$keepLockFile" ]]
    then iconize.skip "rm $packageManagerLock"
    else
      if [[ -f $packageManagerLock ]]
        then iconize.output rm $packageManagerLock
        else iconize.notFound "$packageManagerLock"
      fi
  fi

  if [[ -n "$keepNodeModules" ]]
    then iconize.skip "rm -rf node_modules"
    else
      if [[ -d node_modules ]]
        then spinner rm -rf node_modules
        else iconize.notFound "node_modules"
      fi
  fi

  if [[ -n "$cacheClean" ]]
    then printAndRun 'npm cache clean --force'
    else iconize.skip 'npm cache clean --force'
  fi

  if [[ -n "$skipInstall" ]]
    then iconize.skip "$packageManager install" && return 0
  fi

  printAndRun "$packageManager install" && ${zsb}.tryGitStatus

  alert "🏁 $packageManager install 🏁"
}

compdef "_${zsb}.nonRepeatedListD \
  '-l:Keep lock files' \
  '-n:Keep node modules' \
  '-i:Skip installation' \
  '-f:Force Cache clean' \
  " ${zsb}.nodeReinstall

alias npmreinstall="${zsb}.nodeReinstall npm package-lock.json"

alias yarnreinstall="${zsb}.nodeReinstall yarn yarn.lock"
