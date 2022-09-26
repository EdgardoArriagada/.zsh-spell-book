stash() {
  if [[ -n "$1" ]]
    then git stash push --include-untracked -m "${*}"
    else git stash push --include-untracked
  fi

  gsl
}

alias STASH='toggleCapsLock && stash'

hisIgnore stash

_${zsb}.nocompletion stash

