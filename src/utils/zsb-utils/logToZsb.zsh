logToZsb() (
  local message="$1"
  local folderPath="$2"
  local logFileName="$3"

  if [ -z "$logFileName" ]; then
    logFileName=$(date +%b-%d-%Y.log)
  fi

  local logsPath=${ZSB_DIR}/logs/${folderPath}

  [ ! -d $logsPath ] && mkdir -p $logsPath
  echo "$message" >> "${logsPath}/${logFileName}"
)
