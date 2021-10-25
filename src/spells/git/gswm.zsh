gswm() {
  if git rev-parse --verify master >/dev/null 2>&1; then
    git switch master; return $?
  fi

  git switch main
}

_${zsb}.nocompletion gswm

alias GSWM='toggleCapsLock && gswm'
