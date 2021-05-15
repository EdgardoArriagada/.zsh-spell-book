# This command will find all occurrences of "$1" in all the files under the current directory
# for multiple words, pass the argument surrounded with double quotes, example: dondedice "foo bar"

dondedice() {
  readonly searchTerm=${1:?'Provide a search term.'}
  readonly dir=${2:='.'}

  # You CAN move to any excluded folder to perform a search inside it
  grep -rn --exclude-dir={log,node_modules,build,dist,.cache,coverage,target,tmp,venv} "$searchTerm" "$dir"
  return 0
}
