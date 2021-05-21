gsl() {
  local -r stashList="$(git stash list)"
  [[ -z "$stashList" ]] && ${zsb}.info "Stash list is empty." && return 0
  echo "$stashList" | c -l ruby -
}

alias GSL="toggleCapsLock && gsl"

_${zsb}.nocompletion gsl
