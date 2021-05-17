# golang compile and execute in one go

xgo() {
  local inputFile="$1"
  local extension=$(echo "$inputFile" | rev | cut -d"." -f1 | rev)

  if [[ -z "$inputFile" ]]; then
    echo "${ZSB_ERROR} you have to pass a *.go file as an argument"
    return 1
  fi

  if [[ ! -f "$inputFile" ]]; then
    echo "${ZSB_ERROR} the file ${inputFile} doesn't exists"
    return 1
  fi

  if [[ "$extension" != "go" ]]; then
    echo "${ZSB_ERROR} invalid file extension"
    return 1
  fi

  # outputFile is basically inputFile without the .cpp extension
  local outputFile=$(echo "$inputFile" | rev | cut -d"." -f2 | rev)
  go build "$inputFile" && ./${outputFile}

  return 0
}

# tabulate will autocomplete only files that ends with .cpp
complete -f -X '!*.go' xgo
