${zsb}.searchAndReplace.replaceInFiles() {
  local filesWithMatches=(`rg --files-with-matches ${search}`)

  if (( ${#filesWithMatches} == 0 )); then
    ${zsb}.info "No files found with matches for '`hl ${search}`'"
    return 1
  fi

  if [[ -z "$skipPrompt" ]]; then
    ${zsb}.confirmMenu.withPrompt
  fi

  sd ${search} ${replace} ${filesWithMatches[@]}

  (( $? == 0 )) && ${zsb}.success "`hl ${search}` -> `hl ${replace}`"
}

${zsb}.searchAndReplace.replaceFilesAndFolders() {
  local filesAndFoldersWithMatches=(`fd ${search}`)

  if (( ${#filesAndFoldersWithMatches} == 0 )); then
    ${zsb}.info "No files or folders found with matches for '`hl ${search}`'"
    return 1
  fi

  if [[ -z "$skipPrompt" ]]; then
    ${zsb}.confirmMenu.withPrompt
  fi

  for file in ${filesAndFoldersWithMatches[@]}; do
    if [[ -d ${file} ]]
      then continue
    fi

    local newFile=`sd ${search} ${replace} <<< ${file}`

    mkdir -p `dirname ${newFile}`
    mv ${file} ${newFile}

    (( $? == 0 )) && ${zsb}.success "`hl ${file}` -> `hl ${newFile}`"
  done
}

searchAndReplace() {
  local this=${zsb}.${0}

  zparseopts -D -E -F -- f=filesAndFolders y=skipPrompt || return 1

  local search=${1:?Search string is required}
  local replace=${2:?Replace string is required}

  if [[ -n "$filesAndFolders" ]]
    then ${this}.replaceFilesAndFolders
    else ${this}.replaceInFiles
  fi
}

compdef "_${zsb}.nonRepeatedListD \
  '-f:replace files and dirs' \
  '-c:replace file contents (default)' \
  " searchAndReplace

