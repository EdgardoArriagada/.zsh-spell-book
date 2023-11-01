qaDoctor() (
  local file=${1:?Error: You must provide a file.}

  ${zsb}.doesMatch "$file" "\.qa$" || ${zsb}.throw "File has to be a $(hl ".qa")"

  zmodload zsh/mapfile
  local fileLines=( "${(f)mapfile[$file]}" )

  local Q='Q--'
  local A='A--'

  integer tokenCounter=0
  integer i=0

  for item in $fileLines; do
    (( i++ ))

    if [[ "$item" = "$Q" ]]; then
      (( tokenCounter++ ))
      (( tokenCounter >= 2 )) && ${zsb}.throw "Question duplication in line `hl $i`."
      continue
    fi

    if [[ "$item" = "$A" ]]; then
      (( tokenCounter-- ))
      (( tokenCounter < 0 )) && ${zsb}.throw "Answer duplication in line `hl $i`."
      continue
    fi
  done

  ${zsb}.success "`hl $file` is ok."
)
