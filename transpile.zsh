rm transpiled.zsh
rm temp_transpilation.zsh
# Root global variables
ZSB_DIR=`dirname ${0}`
ZSB_TEMP_DIR=${ZSB_DIR}/src/temp

# Dynamic prefix
: ${zsb:='zsb'}

${zsb}.sourceFiles() for f in ${*}; do cat ${f} >> temp_transpilation.zsh; done

() {
  # Source environment variables
  local envFile=${ZSB_DIR}/.env
  [[ -f ${envFile} ]] && cat ${envFile} >> temp_transpilation.zsh
  cat ${ZSB_DIR}/src/zsh.config.zsh >> temp_transpilation.zsh
  cat ${ZSB_DIR}/src/globalVariables.zsh >> temp_transpilation.zsh

  # Source files in this specific order
  ${zsb}.sourceFiles ${ZSB_DIR}/src/utils/**/*.zsh
  ${zsb}.sourceFiles ${ZSB_DIR}/src/configurations/**/*.zsh
  ${zsb}.sourceFiles ${ZSB_DIR}/src/spells/**/*.zsh

  # Temporal Spells
  local tempSpells=( ${ZSB_DIR}/src/temp/spells**/*.zsh )>/dev/null 2>&1
  [[ -n "$tempSpells" ]] && ${zsb}.sourceFiles ${tempSpells}

  ${zsb}.sourceFiles ${ZSB_DIR}/src/automatic-calls/**/*.zsh
}

sd '( |^)#.*' '' temp_transpilation.zsh # remove comments
sd '\$\{zsb\}' "${zsb}" temp_transpilation.zsh
sd '\$\{ZSB_DIR\}' "${ZSB_DIR}" temp_transpilation.zsh
sd '\$\{ZSB_TEMP_DIR\}' "${ZSB_TEMP_DIR}" temp_transpilation.zsh
sed '/^$/d' temp_transpilation.zsh > transpiled.zsh
rm temp_transpilation.zsh
