gp() {
  {
    declare -A args
    args[--aware]=false
    ${zsb}.switchTrueMatching "${args[@]}" "$@"
    set -- $(${zsb}.clearFlags "${args[@]}" "$@")
  }

  if [[ -z "$1" ]]; then
    echo "${ZSB_ERROR} No branch have been provided"
    return 1
  fi

  local INPUT_REMOTE_BRANCH="$1"
  shift 1

  if ${zsb}.isDefaultBranch "$INPUT_REMOTE_BRANCH" && ! "${args[--aware]}"; then
    echo "${ZSB_ERROR} It's not safe to push to default branch, use ${ZSB_SHL}--aware${ZSB_EHL} flag to do it anyway"
    return 1
  fi

  git push origin "$INPUT_REMOTE_BRANCH" "$@" && copythis " " -s

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
