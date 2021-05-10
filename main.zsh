# Root global variables
ZSB_DIR=$(dirname $0)
ZSB_SNAPSHOT_PATH="$PATH"

() {
  # Source environment variables
  local envFile=${ZSB_DIR}/.env
  if [ -f $envFile ]; then source $envFile; fi

  source ${ZSB_DIR}/src/shell.config.zsh
  source ${ZSB_DIR}/src/zsh.config.zsh

  # Dynamic prefix
  [ -z $zsb ] && zsb="zsb"

  # Source global variables
  source ${ZSB_DIR}/src/globalVariables.zsh

  # Declare source files
  local utilFiles=(${ZSB_DIR}/src/utils/**/*.zsh)>/dev/null 2>&1
  local configurationFiles=(${ZSB_DIR}/src/configurations/**/*.zsh)>/dev/null 2>&1
  local spellPages=(${ZSB_DIR}/src/spells/**/*.zsh)>/dev/null 2>&1
  local automaticCallFiles=(${ZSB_DIR}/src/automatic-calls/**/*.zsh)>/dev/null 2>&1

  # Create a dynamic prefixed function
  __${zsb}.sourceFiles() for file in $*; do source "$file"; done

  # Source files in this specific order
  __${zsb}.sourceFiles $utilFiles
  __${zsb}.sourceFiles $configurationFiles
  __${zsb}.sourceFiles $spellPages
  __${zsb}.sourceFiles $automaticCallFiles

  # Source temporal files (ignored by git)
  local tempDir=${ZSB_DIR}/src/temp
  if [ -d $tempDir ]; then
    local tempFiles=(${tempDir}/**/*.zsh)>/dev/null 2>&1
    __${zsb}.sourceFiles $tempFiles
  fi
}

# Remove dynamic prefixed functions
unfunction -m "__${zsb}.*"
