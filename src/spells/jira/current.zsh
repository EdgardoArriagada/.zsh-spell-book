${zsb}.current.getDir() <<< $ZSB_TICKETS_DIR/$ZSB_PARENT_TICKET

${zsb}_createCurrentDir() {
  ${zsb}.assertJira

  local parentTicketDir=`${zsb}.current.getDir`

  [[ -d $parentTicketDir ]] && return 0

  mkdir -p $parentTicketDir
  print "# $ZSB_CURRENT_LABEL" > $parentTicketDir/NOTES.md
  (builtin cd $parentTicketDir && git init);

  return 0
}

alias cdcurrent="${zsb}_createCurrentDir && cds `${zsb}.current.getDir`"
alias builtincdcurrent="${zsb}_createCurrentDir && builtin cd `${zsb}.current.getDir`"

# We enter in the directory to be ablo to open new shell to interact with git
alias vnotescurrent="${zsb}_createCurrentDir && (builtincdcurrent && nvim `${zsb}.current.getDir`/NOTES.md)"
alias cnotescurrent="${zsb}_createCurrentDir && zsb_cat `${zsb}.current.getDir`/NOTES.md"
alias ncurrent="cdcurrent && (( $ZSB_MACOS )) && open . || nautilus ."

hisIgnore cdcurrent vnotescurrent cnotescurrent ncurrent

pomodorocurrent() {
  ${zsb}.assertIsSet 'ZSB_CURRENT_TICKET'
  local inputTime=${1?Error: missing time}

  ${zsb}.info "`hl $ZSB_CURRENT_TICKET` $ZSB_CURRENT_LABEL"
  pomodoro $inputTime $ZSB_CURRENT_TICKET $ZSB_CURRENT_LABEL
}
