gsl() {
  local -r stashList="$(git stash list)"
  [[ -z "$stashList" ]] && ${zsb}.cancel "Stash list is empty."
  echo "$stashList" | c -p -l ruby -
}

hisIgnore gsl

alias GSL="toggleCapsLock && gsl"

_${zsb}.nocompletion gsl
