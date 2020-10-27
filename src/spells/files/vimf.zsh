# usefull for opening  files found using `grep -rn <input> .`
vimf() (
  local fileInput=$(echo "$1" | cut -d":" -f1)
  local lineNumber=$(echo "$1" | cut -s -d":" -f2)

  main() {
    [ ! -f "$fileInput" ] && echo "${ZSB_ERROR} File not found" && return 1

    if existsLineNumber && ${zsb}_isInteger "$lineNumber"; then
      openFileInNumber
      return 0
    fi

    openFileFound
  }

  existsLineNumber() {
    [ ! -z "$lineNumber" ]
  }

  openFileInNumber() eval "vim +${lineNumber} ${fileInput}"

  openFileFound() eval "vim ${fileInput}"

  main "$@"
)

alias vif="vimf"
alias vf="vimf"
