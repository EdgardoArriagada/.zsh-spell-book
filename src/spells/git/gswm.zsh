_gswm.switch() git switch ${1} >/dev/null 2>&1

gswm() { _gswm.switch 'main' ||  _gswm.switch 'master'; }

hisIgnore gswm

_${zsb}.nocompletion gswm

alias GSWM='toggleCapsLock && gswm'
