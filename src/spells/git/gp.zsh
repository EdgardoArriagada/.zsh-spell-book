gp() {
  {
    declare -A args
    args[--aware]=false
    ${zsb}_switchTrueMatching "${args[@]}" "$@"
    set -- $(${zsb}_clearFlags "${args[@]}" "$@")
  }

  if [ -z "$1" ]; then
    echo "${ZSB_ERROR} No branch have been provided"
    return 1
  fi

  local INPUT_REMOTE_BRANCH="$1"
  shift 1

  if ${zsb}_isDefaultBranch "$INPUT_REMOTE_BRANCH" && ! "${args[--aware]}"; then
    echo "${ZSB_ERROR} It's not safe to push to default branch, use ${ZSB_SHL}--aware${ZSB_EHL} flag to do it anyway"
    return 1
  fi

  git push origin "$INPUT_REMOTE_BRANCH" "$@" && copythis " " -c

  return 0
}

complete -C "isGitRepo && git branch | sed 's/^\*//'" -W '--aware' gp
