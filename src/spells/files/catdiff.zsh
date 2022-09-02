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

    if [[ -z "$file1Lacks" && -z "$file2Lacks" ]]; then
      ${zsb}.info "Both files are equal"
      return 0
    fi

    if [[ -n "$file1Lacks" ]]; then
      ${this}.printFileInfo "$file1"
      echo "$file1Lacks"
    fi

    if [[ -n "$file1Lacks" && -n "$file2Lacks" ]]; then
      hr
      echo " "
    fi

    if [[ -n "$file2Lacks" ]]; then
      ${this}.printFileInfo "$file2"
      echo "$file2Lacks"
    fi

    echo " "
  }

  ${this}.areArgsInvalid() [[ ! -f "$file1" || ! -f "$file2" ]]

  ${this}.throwArgsInvalidException() {
    ${zsb}.throw "Two valid files to compare expected"
  }

  ${this}.printFileInfo() {
    ${zsb}.info "`hl ${1}` lacks of the following lines:"
    echo " "
  }

  ${this}.compareFiles() comm -13 <(sort $1) <(sort $2)

  ${this}.main "$@"
)
