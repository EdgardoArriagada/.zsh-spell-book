catdiff() (
  local this="$0"
  local file1="$1"
  local file2="$2"

  ${this}.main() {
    if ${this}.areArgsInvalid; then
      ${this}.throwArgsInvalidException; return $?
    fi

    local file1Lacks=$(${this}.compareFiles "$file1" "$file2")
    local file2Lacks=$(${this}.compareFiles "$file2" "$file1")

    if [ -z "$file1Lacks" ] && [ -z "$file2Lacks" ]; then
      echo "${ZSB_INFO} Both files are equal."
      return 0
    fi

    if [ ! -z "$file1Lacks" ]; then
      printFileInfo "$file1"
      echo "$file1Lacks"
    fi

    if [ ! -z "$file1Lacks" ] && [ ! -z "$file2Lacks" ]; then
      ${zsb}.printHr
    fi

    if [ ! -z "$file2Lacks" ]; then
      printFileInfo "$file2"
      echo "$file2Lacks"
    fi

    echo " "
  }

  ${this}.areArgsInvalid() $([ ! -f "$file1" ] || [ ! -f "$file2" ])

  ${this}.throwArgsInvalidException() {
    echo "${ZSB_ERROR} Two valid files to compare expected."
    return 1
  }

  ${this}.printFileInfo() {
    echo "${ZSB_INFO} $(hl $1) lacks of the following lines:"
    echo " "
  }

  ${this}.compareFiles() comm -13 <(sort $1) <(sort $2)

  ${this}.main "$@"
)
