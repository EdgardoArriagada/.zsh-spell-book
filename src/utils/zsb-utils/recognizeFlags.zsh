# recognizeFlags
# {
#   local existingFlags=(--one --two)
#   local args=("$@")
#   local -A flags=( ${(z)$(${zsb}.recognizeFlags "args" "existingFlags")} )
#   set -- $(${zsb}.clearFlags "args" "existingFlags")
# }
#
# if "${flags[--one]}"; then...


${zsb}.recognizeFlags() {
  local args=(${(P)1})
  local flags=(${(P)2})
  local usedFlags=( ${flags:*args} )
  local unusedFlags=( ${flags:|usedFlags} )

  local trueFlags=""
  [[ -n "$usedFlags" ]] && trueFlags=$(print -r -- "${(zu)^usedFlags} true")

  local falseFlags=""
  [[ -n "$unusedFlags" ]] && falseFlags=$(print -r -- "${(zu)^unusedFlags} false")

  echo "${trueFlags} ${falseFlags}"
}

${zsb}.clearFlags() {
  local args=(${(P)1})
  local flags=(${(P)2})
  print ${args:|flags}
}
