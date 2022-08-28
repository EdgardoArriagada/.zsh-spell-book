gd() {
  if [[ -z "$1" ]]; then
    git diff && ${zsb}.gitStatus; return $?
  fi

  git diff -- "$@" && ${zsb}.gitStatus
}

hisIgnore gd

compdef "_${zsb}.gitUnrepeat 'red-with-diff'" gd

