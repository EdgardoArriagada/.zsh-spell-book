logToZsb() (
  local message=$1
  local folderPath=$2
  local logFilename=$3

  if [[ -z "$logFilename" ]]; then
    logFilename=$(date +%b-%d-%Y.log)
  fi

  local logsPath=$ZSB_DIR/logs/$folderPath

  [[ ! -d $logsPath ]] && mkdir -p $logsPath
  echo $message >> $logsPath/$logFilename
)
