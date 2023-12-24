${zsb}.searchAndReplace.replaceInFiles() {
  local filesWithMatches=(`rg --files-with-matches $search`)

  if (( $#filesWithMatches == 0 )); then
    [[ -z "$silent" ]] && ${zsb}.info "No files found with matches for '`hl $search`'"
    return 1
  fi

  if [[ -z "$skipPrompt" ]]; then
    ${zsb}.confirmMenu.withPrompt
  fi

  sd $search $replace ${filesWithMatches[@]}

  if [[ -z "$silent" ]]
    then (( $? == 0 )) && ${zsb}.success "`hl $search` -> `hl $replace`"
  fi
}

${zsb}.searchAndReplace.replaceFilesAndFolders() {
  local filesAndFoldersWithMatches=(`fd $search`)

  if (( $#filesAndFoldersWithMatches == 0 )); then
    [[ -z "$silent" ]] && ${zsb}.info "No files or folders found with matches for '`hl $search`'"
    return 1
  fi

  if [[ -z "$skipPrompt" ]]; then
    ${zsb}.confirmMenu.withPrompt
  fi

  for file in ${filesAndFoldersWithMatches[@]}; do
    if [[ -d $file ]]
      then continue
    fi

    local newFile=`sd $search $replace <<< $file`

    mkdir -p `dirname $newFile`
    mv $file $newFile

    if [[ -z "$silent" ]]
      then (( $? == 0 )) && ${zsb}.success "`hl $file` -> `hl $newFile`"
    fi
  done
}

searchAndReplace() {
  local this=${zsb}.${0}

  zparseopts -D -E -F -- f=filesAndFolders y=skipPrompt s=silent || return 1

  local search=${1:?Error: Search string is required.}
  local replace=${2:?Error: Replace string is required.}

  if [[ -n "$filesAndFolders" ]]
    then ${this}.replaceFilesAndFolders
    else ${this}.replaceInFiles
  fi
}

compdef "_${zsb}.nonRepeatedListD \
  '-f:replace files and dirs' \
  '-c:replace file contents (default)' \
  " searchAndReplace

