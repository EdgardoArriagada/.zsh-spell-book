catdiff() (
  local file1="$1"
  local file2="$2"

  main() {
    if areArgsInvalid; then
      throwArgsInvalidException; return $?
    fi

    local file1Lacks=$(compareFiles "$file1" "$file2")
    local file2Lacks=$(compareFiles "$file2" "$file1")

    if [ -z "$file1Lacks" ] && [ -z "$file2Lacks" ]; then
      echo "${ZSB_INFO} Both files are equal."
      return 0
    fi

    if [ ! -z "$file1Lacks" ]; then
      printFileInfo "$file1"
      echo "$file1Lacks"
    fi

    if [ ! -z "$file1Lacks" ] && [ ! -z "$file2Lacks" ]; then
      ${zsb}_printHr
    fi

    if [ ! -z "$file2Lacks" ]; then
      printFileInfo "$file2"
      echo "$file2Lacks"
    fi

    echo " "
  }

  areArgsInvalid() $([ ! -f "$file1" ] || [ ! -f "$file2" ])

  throwArgsInvalidException() {
    echo "${ZSB_ERROR} Two valid files to compare expected."
    return 1
  }

  printFileInfo() {
    echo "${ZSB_INFO} $(hl $1) lacks of the following lines:"
    echo " "
  }

  compareFiles() comm -13 <(sort $1) <(sort $2)

  main "$@"
)
