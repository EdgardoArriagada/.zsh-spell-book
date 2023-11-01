catdiff() (
  local this=$0
  local file1=$1
  local file2=$2

  ${this}.assertValidArgs() {
    if [[ ! -f $file1 || ! -f $file2 ]]
      then ${zsb}.throw "Two valid files to compare expected"
    fi
  }

  ${this}.printFileInfo() {
    ${zsb}.info "`hl $1` lacks of the following lines:\n"
  }

  ${this}.compareFiles() comm -13 <(sort $1) <(sort $2)

  { # main
    ${this}.assertValidArgs

    local file1Lacks=`${this}.compareFiles $file1 $file2`
    local file2Lacks=`${this}.compareFiles $file2 $file1`

    if [[ -z "$file1Lacks" && -z "$file2Lacks" ]]; then
      ${zsb}.info "Both files are equal"
      return 0
    fi

    if [[ -n "$file1Lacks" ]]; then
      ${this}.printFileInfo $file1
      <<< ${file1Lacks}
    fi

    if [[ -n "$file1Lacks" && -n "$file2Lacks" ]]; then
      hr
      <<< ''
    fi

    if [[ -n "$file2Lacks" ]]; then
      ${this}.printFileInfo $file2
      <<< $file2Lacks
    fi

    <<< ''
  }
)
