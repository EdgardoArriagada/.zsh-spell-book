searchAndReplace() {
  local search=${1:?Search string is required}
  local replace=${2:?Replace string is required}

  local filesWithMatches=(`rg --files-with-matches ${search}`)

  if (( ${#filesWithMatches[@]} == 0 )); then
    ${zsb}.throw "No files found with matches for '`hl ${search}`'"
    return 1
  fi

  ${zsb}.confirmMenu.withPrompt

  sd ${search} ${replace} ${filesWithMatches[@]}

  (( $? == 0 )) && ${zsb}.success "`hl ${search}` -> `hl ${replace}`"
}


