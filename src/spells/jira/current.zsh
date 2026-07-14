${zsb}_createCurrentDir() {
  ${zsb}.assertJira

  [[ -d $ZSB_CURRENT_TICKET_DIR ]] && return 0

  mkdir -p $ZSB_CURRENT_TICKET_DIR
  print "# $ZSB_CURRENT_LABEL" > $ZSB_CURRENT_TICKET_DIR/NOTES.md
  (builtin cd $ZSB_PARENT_TICKET_DIR && git init);

  return 0
}

alias cdcurrent="${zsb}_createCurrentDir && cds $ZSB_CURRENT_TICKET_DIR"
alias cinit='cdcurrent'

alias vrepos="nvim $ZSB_CURRENT_TICKET_DIR/.repos"

cdc() {
  local repos=$ZSB_CURRENT_TICKET_DIR/.repos
  [[ ! -f $repos ]] && ${zsb}.throw "`hl $repos` not set"

  local -a entries
  entries=("${(@f)$(grep -v '^[[:space:]]*$' $repos)}")

  local selection
  if (( ${#entries} == 1 ))
    then selection=$entries[1]
    else selection=$(print -l ${entries[@]} | fzf)
  fi

  if [[ -z "$selection" ]]
    then ${zsb}.cancel 'no selection'
  elif [[ -d ${~selection} ]]
    then cd ${~selection}
  else
    ${zsb}.throw "`hl $selection` is not a directory"
  fi
}

hisIgnore cdcurrent vnotescurrent cnotescurrent ncurrent cinit cdc vrepos


pomodorocurrent() {
  ${zsb}.assertIsSet 'ZSB_CURRENT_TICKET'
  local inputTime=${1?Error: missing time}

  ${zsb}.info "`hl $ZSB_CURRENT_TICKET` $ZSB_CURRENT_LABEL"
  pomodoro $inputTime $ZSB_CURRENT_TICKET $ZSB_CURRENT_LABEL
}
