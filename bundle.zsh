[[ -f transpiled.zsh ]] && rm transpiled.zsh
[[ -f temp_transpilation.zsh ]] && rm temp_transpilation.zsh
# Root global variables
ZSB_DIR=`pwd`

# Dynamic prefix
: ${zsb:='zsb'}
print "${zsb}.sourceFiles() for f in \$*; do source \$f; done" >> temp_transpilation.zsh

${zsb}.transpileFiles() for f in $*; do cat $f >> temp_transpilation.zsh; done

() {
  # Source environment variables
  local envFile=$ZSB_DIR/.env
  [[ -f $envFile ]] && cat $envFile >> temp_transpilation.zsh
  cat $ZSB_DIR/src/zsh.config.zsh >> temp_transpilation.zsh
  cat $ZSB_DIR/src/globalVariables.zsh >> temp_transpilation.zsh

  # Source files in this specific order
  ${zsb}.transpileFiles $ZSB_DIR/src/utils/**/*.zsh
  ${zsb}.transpileFiles $ZSB_DIR/src/configurations/**/*.zsh
  ${zsb}.transpileFiles $ZSB_DIR/src/spells/**/*.zsh

  # Temporal Spells
  local tempSpells=( $ZSB_DIR/src/temp/spells**/*.zsh )>/dev/null 2>&1
  [[ -n "$tempSpells" ]] && ${zsb}.transpileFiles $tempSpells

  ${zsb}.transpileFiles $ZSB_DIR/src/automatic-calls/**/*.zsh
}

ZSB_DIR=`print -P %~` # this works in transpilation only
ZSB_TEMP_DIR=$ZSB_DIR/src/temp

sd '( |^)#.*' '' temp_transpilation.zsh # remove comments
sd '\$\{zsb\}' "$zsb" temp_transpilation.zsh
sd '\$ZSB_DIR' "$ZSB_DIR" temp_transpilation.zsh
sd '\$ZSB_TEMP_DIR' "$ZSB_TEMP_DIR" temp_transpilation.zsh
sd '^ *' '' temp_transpilation.zsh # remove leading spaces
sed '/^$/d' temp_transpilation.zsh > result.zsh # remove empty lines
sd '\\\n' '' result.zsh # remove bash line breaks
rm temp_transpilation.zsh

print "ℨ𝔰𝔥 𝔖𝔭𝔢𝔩𝔩𝔟𝔬𝔬𝔨 bundled!!"
