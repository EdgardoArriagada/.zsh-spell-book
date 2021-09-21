gp() {
  {
    local existingFlags=(--aware)
    local args=("$@")
    local -A flags=( ${(z)$(${zsb}.recognizeFlags "args" "existingFlags")} )
    set -- $(${zsb}.clearFlags "args" "existingFlags")
  }

  local inputRemoteBranch="${1:?'You must provede a branch'}"
  shift 1

  if ${zsb}.isDefaultBranch "$inputRemoteBranch" && [[ -z "${flags[--aware]}" ]]; then
    ${zsb}.throw "It's not safe to push to default branch, use $(hl --aware) flag to do it anyway."
  fi

  git push origin "$inputRemoteBranch" "$@"

  return 0
}


_${zsb}.gp() {
  [[ "$CURRENT" -gt 3 ]] && return 0

  local compList
  case $CURRENT in
    2)
      compList=( $(${zsb}.gitBranches 'current') )
      ;;
    3)
      local -r firstItemUsed="${words[2]}"
      if ${zsb}.isDefaultBranch "$firstItemUsed"; then
        compList=( --aware )
      fi
      ;;
  esac

  _describe 'command' compList
}

compdef _${zsb}.gp gp
