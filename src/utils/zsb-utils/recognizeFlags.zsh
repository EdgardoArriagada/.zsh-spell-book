# recognizeFlags
# {
#   local existingFlags=(--one --two)
#   local args=("$@")
#   local -A flags=( ${(z)$(${zsb}.recognizeFlags "args" "existingFlags")} )
#   set -- $(${zsb}.clearFlags "args" "existingFlags")
# }
#
# if [[ -n "${flags[--aware]}" ]]; then...


${zsb}.recognizeFlags() {
  local args=(${(P)1})
  local flags=(${(P)2})
  local usedFlags=( ${flags:*args} )

  [[ -z "$usedFlags" ]] && return 0

  print -r -- "${(uz)^usedFlags} true"
}

${zsb}.clearFlags() {
  local args=(${(P)1})
  local flags=(${(P)2})
  print ${args:|flags}
}
