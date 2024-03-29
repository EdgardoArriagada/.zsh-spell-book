${zsb}.gitSwitchDefault() {
  local branch=`git branch --format "%(refname:short)" --all \
    | rg -m1 "^(origin/)?${1}$" \
    | sd '^origin/' '' \
  `

  if [[ -n "$branch" ]]
    then git switch $branch
    else ${zsb}.throw "No branch found for $1"
  fi
}

alias gswm="${zsb}.gitSwitchDefault '$ZSB_GIT_MAIN_BRANCHES'"
alias gswd="${zsb}.gitSwitchDefault '$ZSB_GIT_DEVELOP_BRANCHES'"

hisIgnore gswm gswd

_${zsb}.nocompletion gswm
_${zsb}.nocompletion gswd

alias GSWM='toggleCapsLock && gswm'
alias GSWD='toggleCapsLock && gswd'
