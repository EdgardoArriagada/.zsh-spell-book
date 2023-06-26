${zsb}.safeLink() {
  local localFile=${ZSB_DIR}/src/dotFiles/${1}
  local destinationFile=${2}

  # Check if setup is already done
  [[ "$(readlink ${destinationFile})" = "$localFile" ]] &&
    ${zsb}.cancel "$(hl ${destinationFile}) setup is already done."

  # Check if destination file is already busy
  [[ -f $destinationFile ]] &&
    ${zsb}.throw "$(hl ${destinationFile}) is already busy. Please back up it manually before proceeding"

  # Try to make destination parent dir if it does not exists
  mkdir -p $(dirname $destinationFile)

  # perform the symlink
  ln -s ${localFile} ${destinationFile} && \
    ${zsb}.success "$(hl ${destinationFile}) file linked"
}

