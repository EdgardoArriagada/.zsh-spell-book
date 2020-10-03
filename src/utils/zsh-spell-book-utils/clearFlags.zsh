# use this function after ${zsb}_switchTrueMatching
# to clear flags from arguments like:
# {
# 	declare -A args
# 	args[--aware]=false
# 	${zsb}_switchTrueMatching "${args[@]}" "$@"
# 	set -- $(${zsb}_clearFlags "${args[@]}" "$@")
# }
# in the previous case, if the flag
# --aware is given, will be clean from $@

${zsb}_clearFlags() {
  local -n args=$1 >/dev/null 2>&1

  # Remove array of args from $@
  shift ${#args[@]}

  for param in "$@"; do
    shift
    [[ "${args[$param]}" == true ]] && continue
    set -- "$@" "$param" # return shifted
  done

  echo "$@"
}
