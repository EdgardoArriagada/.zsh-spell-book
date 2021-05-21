apply() {
  local -r stashList="$(git stash list)"
  [[ -z "$stashList" ]] && ${zsb}.info "Stash list is empty." && return 0

  [[ ! -z "$1" ]] && git stash apply "$1" && return 0

  gsl
  ${zsb}.prompt "Enter stash $(hl number)."
  local choosenStash
  read choosenStash
  git stash apply "$choosenStash"
}

_${zsb}.nocompletion apply
