${zsb}.gitSwitch() {
  local branchReg="^${1}$"

  for b in `${zsb}.gitBranches`; do
    if [[ "$b" =~ "$branchReg" ]]
      then git switch ${b}; return
    fi
  done
}

alias gswm="${zsb}.gitSwitch '${ZSB_GIT_CRITICAL_BRANCHES}'"
alias gswd="${zsb}.gitSwitch '${ZSB_GIT_DEVELOP_BRANCHES}'"

hisIgnore gswm gswd

_${zsb}.nocompletion gswm
_${zsb}.nocompletion gswd

alias GSWM='toggleCapsLock && gswm'
alias GSWD='toggleCapsLock && gswd'
