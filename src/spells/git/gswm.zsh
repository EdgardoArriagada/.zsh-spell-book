gswm() {
  local this=$0

  for b in `${zsb}.gitBranches`; do
    if [[ "$b" == "main" || "$b" == "master" ]]
      then git switch ${b}; return
    fi
  done
}

hisIgnore gswm

_${zsb}.nocompletion gswm

alias GSWM='toggleCapsLock && gswm'
