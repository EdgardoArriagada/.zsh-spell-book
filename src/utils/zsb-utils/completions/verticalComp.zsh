_${zsb}.verticalComp() {
  local inputComp="${(P)1}"
  local formattedComp=( $(print -rl -- "${(z)^inputComp}:") )
  _describe 'command' formattedComp
}
