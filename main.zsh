# Root global variables
ZSB_DIR=$(dirname $0)
ZSB_TEMP_DIR=${ZSB_DIR}/src/temp
ZSB_SNAPSHOT_PATH="$PATH"

# Dynamic prefix
: ${zsb:='zsb'}

${zsb}.sourceFiles() for f in ${*}; do source ${f}; done

() {
  # Source environment variables
  local -r envFile=${ZSB_DIR}/.env
  [[ -f ${envFile} ]] && source ${envFile}
  source ${ZSB_DIR}/src/zsh.config.zsh
  source ${ZSB_DIR}/src/globalVariables.zsh

  # Declare source files
  local utilFiles=( ${ZSB_DIR}/src/utils/**/*.zsh )
  local configurationFiles=( ${ZSB_DIR}/src/configurations/**/*.zsh )
  local spellPages=( ${ZSB_DIR}/src/spells/**/*.zsh )
  local automaticCallFiles=( ${ZSB_DIR}/src/automatic-calls/**/*.zsh )

  # Source files in this specific order
  ${zsb}.sourceFiles ${utilFiles}
  ${zsb}.sourceFiles ${configurationFiles}
  ${zsb}.sourceFiles ${spellPages}

  # Temporal Spells
  local tempSpells=( ${ZSB_DIR}/src/temp/spells**/*.zsh )>/dev/null 2>&1
  [[ -n "$tempSpells" ]] && ${zsb}.sourceFiles ${tempSpells}

  ${zsb}.sourceFiles ${automaticCallFiles}
}

# Remove dynamic prefixed functions that start with fouble underscore
unfunction -m "__${zsb}.*"

