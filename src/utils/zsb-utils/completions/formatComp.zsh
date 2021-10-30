${zsb}.formatComp() {
  local inputComp="${(P)1}"
  print -rl -- "${(z)^inputComp}:"
}
