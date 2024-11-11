tmm() {
  local newName=${1:=$(get_repo_name)}

  tmux rename $1
}

_${zsb}.nocompletion tmm

hisIgnore tmm

alias TMM="toggleCapsLock && tmm"

tmmw() {
  local newName=${1:=$(get_repo_name)}

  tmux rename-window "  $1"
}

_${zsb}.nocompletion tmmw

hisIgnore tmmw

alias TMMW="toggleCapsLock && tmmw"
