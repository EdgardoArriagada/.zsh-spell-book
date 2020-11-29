# example usage:

# funcTest () {
#   {
#     declare -A args
#     args[--aware]=false
#     args[--signed]=false
#     ${zsb}.switchTrueMatching "${args[@]}" "$@"
#   }

#   if "${args[--aware]}" ; then
# 		echo "yes aware!"
# 	fi

#   if "${args[--signed]}" ; then
# 		echo "yes signed!"
# 	fi
# }

# call it like
# $ funcTest --signed
# check it amazing power

${zsb}.switchTrueMatching() {
  local -n args=$1 >/dev/null 2>&1

  # Remove array of flags from local args
  shift ${#args[@]}

  for cmd in "$@"; do
    if [ "${args[$cmd]}" = false ]; then
      args[$cmd]=true
    fi
  done

  declare -p args >/dev/null
}
