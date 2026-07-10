${zsb}_createCurrentDir() {
  ${zsb}.assertJira

  [[ -d $ZSB_CURRENT_TICKET_DIR ]] && return 0

  mkdir -p $ZSB_CURRENT_TICKET_DIR
  print "# $ZSB_CURRENT_LABEL" > $ZSB_CURRENT_TICKET_DIR/NOTES.md
  (builtin cd `dirname $ZSB_CURRENT_TICKET_DIR` && git init);

  return 0
}

alias cdcurrent="${zsb}_createCurrentDir && cds $ZSB_CURRENT_TICKET_DIR"

# We enter in the directory to be ablo to open new shell to interact with git
alias vnotescurrent="${zsb}_createCurrentDir && nvim $ZSB_CURRENT_TICKET_DIR/NOTES.md"
alias cnotescurrent="${zsb}_createCurrentDir && zsb_cat $ZSB_CURRENT_TICKET_DIR/NOTES.md"
alias ncurrent="cdcurrent && (( $ZSB_MACOS )) && open . || nautilus ."

alias cdc='cdcurrent'

hisIgnore cdcurrent vnotescurrent cnotescurrent ncurrent cdc

pomodorocurrent() {
  ${zsb}.assertIsSet 'ZSB_CURRENT_TICKET'
  local inputTime=${1?Error: missing time}

  ${zsb}.info "`hl $ZSB_CURRENT_TICKET` $ZSB_CURRENT_LABEL"
  pomodoro $inputTime $ZSB_CURRENT_TICKET $ZSB_CURRENT_LABEL
}
