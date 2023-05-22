${zsb}.gitSwitch() {
  local branch=`${zsb}.gitBranches | rg "^${1}$"`

  if [[ -n "${branch}" ]]
    then git switch ${branch}
    else ${zsb}.throw "No branch found for ${1}"
  fi
}

alias gswm="${zsb}.gitSwitch '${ZSB_GIT_MAIN_BRANCHES}'"
alias gswd="${zsb}.gitSwitch '${ZSB_GIT_DEVELOP_BRANCHES}'"

hisIgnore gswm gswd

_${zsb}.nocompletion gswm
_${zsb}.nocompletion gswd

alias GSWM='toggleCapsLock && gswm'
alias GSWD='toggleCapsLock && gswd'
