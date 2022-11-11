drop() {
  ${zsb}.validateStashList

  ${zsb}.confirmMenu.warning "You are about to drop `hl 'the last'` stash entry."

  git stash drop

  # show gsl only if stash list is not empty
  # to avoid the cancel message
  ! ${zsb}.isEmptyStashList && gsl
}

alias STASH='toggleCapsLock && drop'

hisIgnore drop

_${zsb}.nocompletion drop

