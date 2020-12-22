gp() {
  {
    declare -A args
    args[--aware]=false
    ${zsb}.switchTrueMatching "${args[@]}" "$@"
    set -- $(${zsb}.clearFlags "${args[@]}" "$@")
  }

  if [ -z "$1" ]; then
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
  case $COMP_CWORD in
    1)
      local currentCompletion="${COMP_WORDS[COMP_CWORD]}"
      if [ -z "$currentCompletion" ]; then
        COMPREPLY=( $(compgen -C "${zsb}.gitBranches 'current'") )
      else
        COMPREPLY=( $(compgen -C "${zsb}.gitBranches") )
      fi
      ;;
    2)
      local firstItemUsed="${COMP_WORDS[1]}"
      if ${zsb}.isDefaultBranch "$firstItemUsed"; then
        COMPREPLY=( $(compgen -W "--aware") )
      fi
      ;;
  esac
}

complete -F _${zsb}.gp gp
