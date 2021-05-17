# usefull for opening  files found using `grep -rn <input> .`
vimf() (
  local this="$0"
  local fileInput=$(echo "$1" | cut -d":" -f1)
  local lineNumber=$(echo "$1" | cut -s -d":" -f2)

  ${this}.main() {
    [[ ! -f "$fileInput" ]] && echo "${ZSB_ERROR} File not found" && return 1

    if ${this}.existsLineNumber && ${zsb}.isInteger "$lineNumber"; then
      ${this}.openFileInNumber
      return 0
    fi

    ${this}.openFileFound
  }

  ${this}.existsLineNumber() {
    [[ ! -z "$lineNumber" ]]
  }

  ${this}.openFileInNumber() eval "vim +${lineNumber} ${fileInput}"

  ${this}.openFileFound() eval "vim ${fileInput}"

  ${this}.main "$@"
)

alias vif="vimf"
alias vf="vimf"
