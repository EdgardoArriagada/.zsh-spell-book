gp() {
  zparseopts -D -E -F -- -aware=aware || return 1

  local inputRemoteBranch="${1:?'You must provede a branch'}"
  shift 1

  if ${zsb}.isDefaultBranch "$inputRemoteBranch" && [[ -z "$aware" ]]; then
    ${zsb}.throw "It's not safe to push to default branch, use $(hl --aware) flag to do it anyway."
  fi

  ${zsb}.info "Pushing to `hl ${inputRemoteBranch}`"
  git push origin "$inputRemoteBranch" "$@"
}


_${zsb}.gp() {
  (( $CURRENT > 3 )) && return 0

  local compList
  case $CURRENT in
    2)
      compList=( `${zsb}.gitBranches 'current'` )
      ;;
    3)
      local -r firstItemUsed=${words[2]}
      if ${zsb}.isDefaultBranch ${firstItemUsed}; then
        compList=( --aware )
      fi
      ;;
  esac

  _describe 'command' compList
}

compdef _${zsb}.gp gp
