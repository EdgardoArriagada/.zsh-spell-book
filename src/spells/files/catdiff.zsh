catdiff() (
  local file_1="$1"
  local file_2="$2"

  main() {
    if areArgsInvalid; then
      throwArgsInvalidException; return $?
    fi

    local file_1_lacks=$(compareFiles "$file_1" "$file_2")
    local file_2_lacks=$(compareFiles "$file_2" "$file_1")

    if [ -z "$file_1_lacks" ] && [ -z "$file_2_lacks" ]; then
      echo "${ZSB_INFO} Both files are equal."
      return 0
    fi

    if [ ! -z "$file_1_lacks" ]; then
      printFileInfo "$file_1"
      echo "$file_1_lacks"
    fi

    if [ ! -z "$file_1_lacks" ] && [ ! -z "$file_2_lacks" ]; then
      ${zsb}_printHr
    fi

    if [ ! -z "$file_2_lacks" ]; then
      printFileInfo "$file_2"
      echo "$file_2_lacks"
    fi

    echo " "
  }

  areArgsInvalid() return $([ ! -f "$file_1" ] || [ ! -f "$file_2" ])

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
