cdtemp() {
  local tempDir=${ZSB_DIR}/src/temp

  [[ ! -d ${tempDir} ]] && mkdir ${tempDir}

  cds ${tempDir}
}
