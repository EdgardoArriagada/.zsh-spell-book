${zsb}.nodeReinstall() {
  local packageManager=$1
  local packageManagerLock=$2

  zparseopts -D -E -F -- l=keepLockFile n=keepNodeModules i=skipInstall c=cacheClean || return 1

  if [[ ! -f $packageManagerLock ]]
    then iconize.notFound "rm $packageManagerLock"
    else
      if [[ -z "$keepLockFile" ]]
        then iconize.output rm $packageManagerLock
        else iconize.skip "rm $packageManagerLock"
      fi
  fi

  if [[ ! -d node_modules ]]
    then iconize.notFound "rm -rf node_modules"
    else
      if [[ -z "$keepNodeModules" ]]
        then spinner rm -rf node_modules
        else iconize.skip "rm -rf node_modules"
      fi
  fi


  if [[ -n "$cacheClean" ]]
    then printAndRun 'npm cache clean --force'
    else messageSkip 'npm cache clean --force'
  fi

  if [[ -n "$skipInstall" ]]
    then messageSkip "$packageManager install" && return 0
  fi

  printAndRun "$packageManager install" && ${zsb}.isGitRepo && ${zsb}.gitStatus

  alert "🏁 $packageManager install 🏁"
}

compdef "_${zsb}.nonRepeatedListD '-l:Keep lock files' '-n:Keep node modules' '-i:Skip installation' '-c:Cache clean'" ${zsb}.nodeReinstall

alias npmreinstall="${zsb}.nodeReinstall npm package-lock.json"

alias yarnreinstall="${zsb}.nodeReinstall yarn yarn.lock"
