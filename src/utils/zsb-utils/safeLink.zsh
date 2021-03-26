${zsb}.safeLink() {
  local localConfigFile=${ZSB_DIR}/src/configurations/${1}
  local genericConfigFile=${2}

  if [ "$(readlink -f ${genericConfigFile})" = "$localConfigFile" ]; then
    echo "${ZSB_INFO} $(hl ${genericConfigFile}) setup is already done."
    return 0
  fi

  if [ -f $genericConfigFile ]; then
    echo "${ZSB_ERROR} $(hl ${genericConfigFile}) is already busy. Please back up it manually before proceeding"
    return 1
  fi

  ln -s ${localConfigFile} ${genericConfigFile} && \
    echo "${ZSB_SUCCESS} $(hl ${genericConfigFile}) file linked"
}

