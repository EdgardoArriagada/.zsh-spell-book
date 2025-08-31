# Root global variables
ZSB_DIR=`dirname $0`
ZSB_TEMP_DIR=$ZSB_DIR/src/temp

# Dynamic prefix
: ${zsb:='zsb'}

${zsb}.sourceFiles() for f in $*; do source $f; done

() {
  # Source environment variables
  local envFile=$ZSB_DIR/.env
  [[ -f $envFile ]] && source $envFile
  source $ZSB_DIR/src/zsh.config.zsh
  source $ZSB_DIR/src/globalVariables.zsh

  # Source files in this specific order
  ${zsb}.sourceFiles $ZSB_DIR/src/utils/**/*.zsh
  ${zsb}.sourceFiles $ZSB_DIR/src/configurations/**/*.zsh
  ${zsb}.sourceFiles $ZSB_DIR/src/spells/**/*.zsh

  # Temporal Spells
  local tempSpells=( $ZSB_DIR/src/temp/spells**/*.zsh )>/dev/null 2>&1
  [[ -n "$tempSpells" ]] && ${zsb}.sourceFiles $tempSpells

  ${zsb}.sourceFiles $ZSB_DIR/src/automatic-calls/**/*.zsh
}

print "               🕯️ 📜🪶\n"

