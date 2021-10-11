# usefull for opening  files found using `grep -rn <input> .`
# vim found
vf() (
  local this="$0"
  local fileInput="${1%%:*}"
  local lineNumber="$(echo "$1" | cut -s -d":" -f2)"

  ${this}.main() {
    [[ ! -f "$fileInput" ]] && ${zsb}.throw "File not found"

    if ${this}.existsLineNumber && ${zsb}.isInteger "$lineNumber"; then
      ${this}.openFileInNumber
      return 0
    fi

    ${this}.openFileFound
  }

  ${this}.existsLineNumber() [[ -n "$lineNumber" ]]

  ${this}.openFileInNumber() eval "nnvim +${lineNumber} ${fileInput}"

  ${this}.openFileFound() eval "nnvim ${fileInput}"

  ${this}.main "$@"
)

