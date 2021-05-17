searchFileAndOpen() {
  local programToOpenFile="$1"
  local fileToSearch="$2"

  if [[ -z "$fileToSearch" ]]; then
    echo "${ZSB_ERROR} You must specify a file name"
    return 1
  fi

  local sanitizedInput=$(echo "$fileToSearch" | rev | cut -d"/" -f1 | rev)
  local parentDir="$(dirname $fileToSearch)"

  local FOUND_FILE=$(fastFind "$parentDir" "$sanitizedInput" "f")

  if [[ ! -f "$FOUND_FILE" ]]; then
    local tailMessage
    if [[ "$parentDir" != "." ]]; then tailMessage="in $(hl "$parentDir") directory"; fi

    echo "${ZSB_INFO} the file $(hl "$sanitizedInput") was not found ${tailMessage}"
    return 1
  fi

  echo "${ZSB_INFO} File found: $(hl "$FOUND_FILE")"

  eval "${programToOpenFile} ${FOUND_FILE}"
}

alias catg="searchFileAndOpen cat"
alias vimg="searchFileAndOpen vim"
alias vg="vimg"
