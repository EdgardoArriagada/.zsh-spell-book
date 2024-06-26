rename() {
  local inputFile=${1?:Error: You must give a input.}
  local newName=${2?:Error: You must give a new name.}
  local newFile=${newName}.${1:e}

  mv $inputFile $newFile
  ${zsb}.success "`hl $inputFile` -> `hl $newFile`"
}
