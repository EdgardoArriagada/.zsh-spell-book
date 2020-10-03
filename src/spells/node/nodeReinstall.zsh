nodeReinstall() {
  local packageManager=$1
  local packageManagerLock=$2

  {
    local flags=$3
    local GROUP_FLAGS='lni'

    ! areFlagsInGroup "$flags" "$GROUP_FLAGS" && return 1
  }

  if [[ ! "$flags" == *"l"* ]] && [ -f ${packageManagerLock} ]; then
    rm ${packageManagerLock} && echo "${ZSB_INFO} $(hl ${packageManagerLock}) deleted"
  fi

  if [[ ! "$flags" == *"n"* ]] && [ -d node_modules ]; then
    rm -rf node_modules && echo "${ZSB_INFO} $(hl 'node_modules') deleted"
  fi

  if [[ "$flags" == *"i"* ]]; then
    echo "${ZSB_INFO} installation cancelled"
    return 0
  fi

  # BUG: npm always tries to acess to zsb_lazy_load inside a bash function
  # eval is a workaround
  eval "${packageManager} install" && isGitRepo && gitStatus
}

complete -W "-l -n -i -lni" nodeReinstall

alias npmreinstall="nodeReinstall npm package-lock.json"

alias yarnreinstall="nodeReinstall yarn yarn.lock"
