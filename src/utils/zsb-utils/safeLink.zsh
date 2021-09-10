${zsb}.safeLink() {
  local -r localFile=${ZSB_DIR}/src/dotFiles/${1}
  local -r destinationFile=${2}

  # Check if setup is already done
  if [[ "$(readlink -f ${destinationFile})" = "$localFile" ]]; then
    ${zsb}.info "$(hl ${destinationFile}) setup is already done."
    return 0
  fi

  # Check if destination file is already busy
  if [[ -f $destinationFile ]]; then
    ${zsb}.throw "$(hl ${destinationFile}) is already busy. Please back up it manually before proceeding"
  fi

  # Try to make destination parent dir if it does not exists
  mkdir -p $(dirname $destinationFile)

  # perform the symlink
  ln -s ${localFile} ${destinationFile} && \
    ${zsb}.success "$(hl ${destinationFile}) file linked"
}

