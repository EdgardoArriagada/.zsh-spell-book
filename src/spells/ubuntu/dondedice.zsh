# This command will find all occurrences of "$1" in all the files under the current directory
# for multiple words, pass the argument surrounded with double quotes, example: dondedice "foo bar"

dondedice() {
  if [ -z "$1" ]; then
    echo "${ZSB_ERROR} At least one argument expected"
    return 1
  fi

  local dir="$2"

  if [ -z "$dir" ]; then
    dir="."
  fi

  # You CAN move to any excluded folder to perform a search inside it
  grep -rn --exclude-dir={log,node_modules,build,dist,public,.cache,coverage,target,tmp} "$1" "$dir"
  return 0
}
